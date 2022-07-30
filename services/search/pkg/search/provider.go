package search

import (
	"context"
	"fmt"
	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	rpcv1beta1 "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	ctxpkg "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/cs3org/reva/v2/pkg/errtypes"
	sdk "github.com/cs3org/reva/v2/pkg/sdk/common"
	"github.com/cs3org/reva/v2/pkg/storage/utils/walker"
	"github.com/cs3org/reva/v2/pkg/storagespace"
	"github.com/cs3org/reva/v2/pkg/utils"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/search/pkg/content"
	"github.com/owncloud/ocis/v2/services/search/pkg/engine"
	"google.golang.org/grpc/metadata"
	"path/filepath"
	"sort"
	"strings"

	searchmsg "github.com/owncloud/ocis/v2/protogen/gen/ocis/messages/search/v0"
	searchsvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/search/v0"
)

type Provider struct {
	logger    log.Logger
	gateway   gateway.GatewayAPIClient
	engine    engine.Engine
	extractor content.Extractor
	secret    string
}

func NewProvider(gw gateway.GatewayAPIClient, eng engine.Engine, extractor content.Extractor, logger log.Logger, secret string) *Provider {
	return &Provider{
		gateway:   gw,
		engine:    eng,
		secret:    secret,
		logger:    logger,
		extractor: extractor,
	}
}

func (p *Provider) Search(ctx context.Context, req *searchsvc.SearchRequest) (*searchsvc.SearchResponse, error) {
	if req.Query == "" {
		return nil, errtypes.BadRequest("empty query provided")
	}
	p.logger.Debug().Str("query", req.Query).Msg("performing a search")

	listSpacesRes, err := p.gateway.ListStorageSpaces(ctx, &provider.ListStorageSpacesRequest{
		Filters: []*provider.ListStorageSpacesRequest_Filter{
			{
				Type: provider.ListStorageSpacesRequest_Filter_TYPE_SPACE_TYPE,
				Term: &provider.ListStorageSpacesRequest_Filter_SpaceType{SpaceType: "+grant"},
			},
		},
	})
	if err != nil {
		p.logger.Error().Err(err).Msg("failed to list the user's storage spaces")
		return nil, err
	}

	mountpointMap := map[string]string{}
	for _, space := range listSpacesRes.StorageSpaces {
		if space.SpaceType != "mountpoint" {
			continue
		}
		opaqueMap := sdk.DecodeOpaqueMap(space.Opaque)
		grantSpaceId := storagespace.FormatResourceID(provider.ResourceId{
			StorageId: opaqueMap["grantStorageID"],
			SpaceId:   opaqueMap["grantSpaceID"],
			OpaqueId:  opaqueMap["grantOpaqueID"],
		})
		mountpointMap[grantSpaceId] = space.Id.OpaqueId
	}

	matches := MatchArray{}
	total := int32(0)
	for _, space := range listSpacesRes.StorageSpaces {
		searchRootId := &searchmsg.ResourceID{
			StorageId: space.Root.StorageId,
			SpaceId:   space.Root.SpaceId,
			OpaqueId:  space.Root.OpaqueId,
		}
		var mountpointRootId *searchmsg.ResourceID
		mountpointPrefix := ""
		switch space.SpaceType {
		case "mountpoint":
			continue // mountpoint spaces are only "links" to the shared spaces. we have to search the shared "grant" space instead
		case "grant":
			// In case of grant spaces we search the root of the outer space and translate the paths to the according mountpoint
			searchRootId.OpaqueId = space.Root.SpaceId
			mountpointId, ok := mountpointMap[space.Id.OpaqueId]
			if !ok {
				p.logger.Warn().Interface("space", space).Msg("could not find mountpoint space for grant space")
				continue
			}
			gpRes, err := p.gateway.GetPath(ctx, &provider.GetPathRequest{
				ResourceId: space.Root,
			})
			if err != nil {
				p.logger.Error().Err(err).Str("space", space.Id.OpaqueId).Msg("failed to get path for grant space root")
				continue
			}
			if gpRes.Status.Code != rpcv1beta1.Code_CODE_OK {
				p.logger.Error().Interface("status", gpRes.Status).Str("space", space.Id.OpaqueId).Msg("failed to get path for grant space root")
				continue
			}
			mountpointPrefix = utils.MakeRelativePath(gpRes.Path)
			sid, spid, oid, err := storagespace.SplitID(mountpointId)
			if err != nil {
				p.logger.Error().Err(err).Str("space", space.Id.OpaqueId).Str("mountpointId", mountpointId).Msg("invalid mountpoint space id")
				continue
			}
			mountpointRootId = &searchmsg.ResourceID{
				StorageId: sid,
				SpaceId:   spid,
				OpaqueId:  oid,
			}
			p.logger.Debug().Interface("grantSpace", space).Interface("mountpointRootId", mountpointRootId).Msg("searching a grant")
		}

		res, err := p.engine.Search(ctx, &searchsvc.SearchIndexRequest{
			Query: req.Query,
			Ref: &searchmsg.Reference{
				ResourceId: searchRootId,
				Path:       mountpointPrefix,
			},
			PageSize: req.PageSize,
		})
		if err != nil {
			p.logger.Error().Err(err).Str("space", space.Id.OpaqueId).Msg("failed to search the index")
			return nil, err
		}
		p.logger.Debug().Str("space", space.Id.OpaqueId).Int("hits", len(res.Matches)).Msg("space search done")

		total += res.TotalMatches
		for _, match := range res.Matches {
			if mountpointPrefix != "" {
				match.Entity.Ref.Path = utils.MakeRelativePath(strings.TrimPrefix(match.Entity.Ref.Path, mountpointPrefix))
			}
			if mountpointRootId != nil {
				match.Entity.Ref.ResourceId = mountpointRootId
			}
			matches = append(matches, match)
		}
	}

	// compile one sorted list of matches from all spaces and apply the limit if needed
	sort.Sort(matches)
	limit := req.PageSize
	if limit == 0 {
		limit = 200
	}
	if int32(len(matches)) > limit {
		matches = matches[0:limit]
	}

	return &searchsvc.SearchResponse{
		Matches:      matches,
		TotalMatches: total,
	}, nil
}

func (p *Provider) IndexSpace(ctx context.Context, req *searchsvc.IndexSpaceRequest) (*searchsvc.IndexSpaceResponse, error) {
	// get user
	res, err := p.gateway.GetUserByClaim(context.Background(), &user.GetUserByClaimRequest{
		Claim: "username",
		Value: req.UserId,
	})
	if err != nil || res.Status.Code != rpc.Code_CODE_OK {
		fmt.Println("error: Could not get user by userid")
		return nil, err
	}

	// Get auth context
	ownerCtx := ctxpkg.ContextSetUser(context.Background(), res.User)
	authRes, err := p.gateway.Authenticate(ownerCtx, &gateway.AuthenticateRequest{
		Type:         "machine",
		ClientId:     "userid:" + res.User.Id.OpaqueId,
		ClientSecret: p.secret,
	})
	if err != nil || authRes.GetStatus().GetCode() != rpc.Code_CODE_OK {
		return nil, err
	}

	if authRes.GetStatus().GetCode() != rpc.Code_CODE_OK {
		return nil, fmt.Errorf("could not get authenticated context for user")
	}
	ownerCtx = metadata.AppendToOutgoingContext(ownerCtx, ctxpkg.TokenHeader, authRes.Token)

	// Walk the space and index all files
	w := walker.NewWalker(p.gateway)
	rootId, err := storagespace.ParseID(req.SpaceId)
	if err != nil {
		p.logger.Error().Err(err).Msg(err.Error())
		return nil, err
	}
	err = w.Walk(ownerCtx, &rootId, func(wd string, info *provider.ResourceInfo, err error) error {
		if err != nil {
			p.logger.Error().Err(err).Msg("error walking the tree")
		}
		ref := &provider.Reference{
			Path:       utils.MakeRelativePath(filepath.Join(wd, info.Path)),
			ResourceId: &rootId,
		}

		doc, err := p.extractor.Extract(ctx, ref, info)
		if err != nil {
			p.logger.Error().Err(err).Msg("error extracting content")
		}

		ent := engine.Entity{
			ID:       storagespace.FormatResourceID(*info.Id),
			RootID:   storagespace.FormatResourceID(*ref.ResourceId),
			Path:     ref.Path,
			Type:     uint64(info.Type),
			Document: doc,
		}

		err = p.engine.Upsert(ent.ID, ent)
		if err != nil {
			p.logger.Error().Err(err).Msg("error adding resource to the index")
		} else {
			p.logger.Debug().Interface("ref", ref).Msg("added resource to index")
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	logDocCount(p.engine, p.logger)
	return &searchsvc.IndexSpaceResponse{}, nil
}
