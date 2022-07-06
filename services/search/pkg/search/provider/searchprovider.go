package provider

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	rpcv1beta1 "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	ctxpkg "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/cs3org/reva/v2/pkg/errtypes"
	"github.com/cs3org/reva/v2/pkg/events"
	sdk "github.com/cs3org/reva/v2/pkg/sdk/common"
	"github.com/cs3org/reva/v2/pkg/storage/utils/walker"
	"github.com/cs3org/reva/v2/pkg/storagespace"
	"github.com/cs3org/reva/v2/pkg/utils"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/search/pkg/search"
	"google.golang.org/grpc/metadata"

	searchmsg "github.com/owncloud/ocis/v2/protogen/gen/ocis/messages/search/v0"
	searchsvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/search/v0"
)

var ListenEvents = []events.Unmarshaller{
	events.ItemTrashed{},
	events.ItemRestored{},
	events.ItemMoved{},
	events.ContainerCreated{},
	events.FileUploaded{},
	events.FileTouched{},
	events.FileVersionRestored{},
}

type Provider struct {
	logger            log.Logger
	gwClient          gateway.GatewayAPIClient
	indexClient       search.IndexClient
	machineAuthAPIKey string
}

func New(gwClient gateway.GatewayAPIClient, indexClient search.IndexClient, machineAuthAPIKey string, eventsChan <-chan interface{}, logger log.Logger) *Provider {
	p := &Provider{
		gwClient:          gwClient,
		indexClient:       indexClient,
		machineAuthAPIKey: machineAuthAPIKey,
		logger:            logger,
	}

	go func() {
		for {
			ev := <-eventsChan
			go func() {
				time.Sleep(1 * time.Second) // Give some time to let everything settle down before trying to access it when indexing
				p.handleEvent(ev)
			}()
		}
	}()

	return p
}

func (p *Provider) Search(ctx context.Context, req *searchsvc.SearchRequest) (*searchsvc.SearchResponse, error) {
	if req.Query == "" {
		return nil, errtypes.BadRequest("empty query provided")
	}
	p.logger.Debug().Str("query", req.Query).Msg("performing a search")

	listSpacesRes, err := p.gwClient.ListStorageSpaces(ctx, &provider.ListStorageSpacesRequest{
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
			OpaqueId:  opaqueMap["grantOpaqueID"],
		})
		mountpointMap[grantSpaceId] = space.Id.OpaqueId
	}

	matches := []*searchmsg.Match{}
	for _, space := range listSpacesRes.StorageSpaces {
		var mountpointRootId *searchmsg.ResourceID
		mountpointPrefix := ""
		switch space.SpaceType {
		case "mountpoint":
			continue // mountpoint spaces are only "links" to the shared spaces. we have to search the shared "grant" space instead
		case "grant":
			mountpointId, ok := mountpointMap[space.Id.OpaqueId]
			if !ok {
				p.logger.Warn().Interface("space", space).Msg("could not find mountpoint space for grant space")
				continue
			}
			gpRes, err := p.gwClient.GetPath(ctx, &provider.GetPathRequest{
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
			sid, oid, err := storagespace.SplitID(mountpointId)
			if err != nil {
				p.logger.Error().Err(err).Str("space", space.Id.OpaqueId).Str("mountpointId", mountpointId).Msg("invalid mountpoint space id")
				continue
			}
			mountpointRootId = &searchmsg.ResourceID{
				StorageId: sid,
				OpaqueId:  oid,
			}
			p.logger.Debug().Interface("grantSpace", space).Interface("mountpointRootId", mountpointRootId).Msg("searching a grant")
		}

		_, rootStorageID := storagespace.SplitStorageID(space.Root.StorageId)
		res, err := p.indexClient.Search(ctx, &searchsvc.SearchIndexRequest{
			Query: formatQuery(req.Query),
			Ref: &searchmsg.Reference{
				ResourceId: &searchmsg.ResourceID{
					StorageId: space.Root.StorageId,
					OpaqueId:  rootStorageID,
				},
				Path: mountpointPrefix,
			},
			PageSize: req.PageSize,
		})
		if err != nil {
			p.logger.Error().Err(err).Str("space", space.Id.OpaqueId).Msg("failed to search the index")
			return nil, err
		}
		p.logger.Debug().Str("space", space.Id.OpaqueId).Int("hits", len(res.Matches)).Msg("space search done")

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

	return &searchsvc.SearchResponse{
		Matches: matches,
	}, nil
}

func (p *Provider) IndexSpace(ctx context.Context, req *searchsvc.IndexSpaceRequest) (*searchsvc.IndexSpaceResponse, error) {
	// get user
	res, err := p.gwClient.GetUserByClaim(context.Background(), &user.GetUserByClaimRequest{
		Claim: "username",
		Value: req.UserId,
	})
	if err != nil || res.Status.Code != rpc.Code_CODE_OK {
		fmt.Println("error: Could not get user by userid")
		return nil, err
	}

	// Get auth context
	ownerCtx := ctxpkg.ContextSetUser(context.Background(), res.User)
	authRes, err := p.gwClient.Authenticate(ownerCtx, &gateway.AuthenticateRequest{
		Type:         "machine",
		ClientId:     "userid:" + res.User.Id.OpaqueId,
		ClientSecret: p.machineAuthAPIKey,
	})
	if err != nil || authRes.GetStatus().GetCode() != rpc.Code_CODE_OK {
		return nil, err
	}

	if authRes.GetStatus().GetCode() != rpc.Code_CODE_OK {
		return nil, fmt.Errorf("could not get authenticated context for user")
	}
	ownerCtx = metadata.AppendToOutgoingContext(ownerCtx, ctxpkg.TokenHeader, authRes.Token)

	// Walk the space and index all files
	walker := walker.NewWalker(p.gwClient)
	rootId := &provider.ResourceId{StorageId: req.SpaceId, OpaqueId: req.SpaceId}
	err = walker.Walk(ownerCtx, rootId, func(wd string, info *provider.ResourceInfo, err error) error {
		if err != nil {
			p.logger.Error().Err(err).Msg("error walking the tree")
		}
		ref := &provider.Reference{
			Path:       utils.MakeRelativePath(filepath.Join(wd, info.Path)),
			ResourceId: rootId,
		}
		err = p.indexClient.Add(ref, info)
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

	p.logDocCount()
	return &searchsvc.IndexSpaceResponse{}, nil
}

func (p *Provider) logDocCount() {
	c, err := p.indexClient.DocCount()
	if err != nil {
		p.logger.Error().Err(err).Msg("error getting document count from the index")
	}
	p.logger.Debug().Interface("count", c).Msg("new document count")
}

func formatQuery(q string) string {
	query := q
	fields := []string{"RootID", "Path", "ID", "Name", "Size", "Mtime", "MimeType", "Type"}
	for _, field := range fields {
		query = strings.ReplaceAll(query, strings.ToLower(field)+":", field+":")
	}

	if strings.Contains(query, ":") {
		return query // Sophisticated field based search
	}

	// this is a basic filename search
	return "Name:*" + strings.ReplaceAll(strings.ToLower(query), " ", `\ `) + "*"
}
