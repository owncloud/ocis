package provider

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	rpcv1beta1 "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	typesv1beta1 "github.com/cs3org/go-cs3apis/cs3/types/v1beta1"
	ctxpkg "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/cs3org/reva/v2/pkg/errtypes"
	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/cs3org/reva/v2/pkg/storage/utils/walker"
	"github.com/cs3org/reva/v2/pkg/utils"
	"github.com/owncloud/ocis/extensions/search/pkg/search"
	"github.com/owncloud/ocis/ocis-pkg/log"
	"google.golang.org/grpc/metadata"

	searchmsg "github.com/owncloud/ocis/protogen/gen/ocis/messages/search/v0"
	searchsvc "github.com/owncloud/ocis/protogen/gen/ocis/services/search/v0"
)

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
			var ref *provider.Reference
			var owner *user.User
			switch e := ev.(type) {
			case events.ItemTrashed:
				p.logger.Debug().Interface("event", ev).Msg("marking document as deleted")
				err := p.indexClient.Delete(e.ID)
				if err != nil {
					p.logger.Error().Err(err).Interface("Id", e.ID).Msg("failed to remove item from index")
				}
				continue
			case events.ItemRestored:
				p.logger.Debug().Interface("event", ev).Msg("marking document as restored")
				ref = e.Ref
				owner = &user.User{
					Id: e.Executant,
				}

				statRes, err := p.statResource(ref, owner)
				if err != nil {
					p.logger.Error().Err(err).Msg("failed to stat the changed resource")
				}

				switch statRes.Status.Code {
				case rpc.Code_CODE_OK:
					err = p.indexClient.Restore(statRes.Info.Id)
					if err != nil {
						p.logger.Error().Err(err).Msg("failed to restore the changed resource in the index")
					}
				default:
					p.logger.Error().Interface("statRes", statRes).Msg("failed to stat the changed resource")
				}

				continue
			case events.ItemMoved:
				p.logger.Debug().Interface("event", ev).Msg("resource has been moved, updating the document")
				ref = e.Ref
				owner = &user.User{
					Id: e.Executant,
				}

				statRes, err := p.statResource(ref, owner)
				if err != nil {
					p.logger.Error().Err(err).Msg("failed to stat the changed resource")
				}

				switch statRes.Status.Code {
				case rpc.Code_CODE_OK:
					err = p.indexClient.Move(statRes.Info)
					if err != nil {
						p.logger.Error().Err(err).Msg("failed to restore the changed resource in the index")
					}
				default:
					p.logger.Error().Interface("statRes", statRes).Msg("failed to stat the changed resource")
				}

				continue
			case events.ContainerCreated:
				ref = e.Ref
				owner = &user.User{
					Id: e.Executant,
				}
			case events.FileUploaded:
				ref = e.Ref
				owner = &user.User{
					Id: e.Executant,
				}
			case events.FileVersionRestored:
				ref = e.Ref
				owner = &user.User{
					Id: e.Executant,
				}
			default:
				// Not sure what to do here. Skip.
				continue
			}
			p.logger.Debug().Interface("event", ev).Msg("resource has been changed, updating the document")

			statRes, err := p.statResource(ref, owner)
			if err != nil {
				p.logger.Error().Err(err).Msg("failed to stat the changed resource")
			}

			switch statRes.Status.Code {
			case rpc.Code_CODE_OK:
				err = p.indexClient.Add(ref, statRes.Info)
				if err != nil {
					p.logger.Error().Err(err).Msg("error adding updating the resource in the index")
				} else {
					p.logDocCount()
				}
			default:
				p.logger.Error().Interface("statRes", statRes).Msg("failed to stat the changed resource")
			}
		}
	}()

	return p
}

func (p *Provider) statResource(ref *provider.Reference, owner *user.User) (*provider.StatResponse, error) {
	// Get auth
	ownerCtx := ctxpkg.ContextSetUser(context.Background(), owner)
	authRes, err := p.gwClient.Authenticate(ownerCtx, &gateway.AuthenticateRequest{
		Type:         "machine",
		ClientId:     "userid:" + owner.Id.OpaqueId,
		ClientSecret: p.machineAuthAPIKey,
	})
	if err != nil || authRes.GetStatus().GetCode() != rpc.Code_CODE_OK {
		p.logger.Error().Err(err).Interface("authRes", authRes).Msg("error using machine auth")
	}
	ownerCtx = metadata.AppendToOutgoingContext(ownerCtx, ctxpkg.TokenHeader, authRes.Token)

	// Stat changed resource resource
	return p.gwClient.Stat(ownerCtx, &provider.StatRequest{Ref: ref})
}

func (p *Provider) logDocCount() {
	c, err := p.indexClient.DocCount()
	if err != nil {
		p.logger.Error().Err(err).Msg("error getting document count from the index")
	}
	p.logger.Debug().Interface("count", c).Msg("new document count")
}

func (p *Provider) Search(ctx context.Context, req *searchsvc.SearchRequest) (*searchsvc.SearchResponse, error) {
	if req.Query == "" {
		return nil, errtypes.PreconditionFailed("empty query provided")
	}

	listSpacesRes, err := p.gwClient.ListStorageSpaces(ctx, &provider.ListStorageSpacesRequest{
		Opaque: &typesv1beta1.Opaque{Map: map[string]*typesv1beta1.OpaqueEntry{
			"path": {
				Decoder: "plain",
				Value:   []byte("/"),
			},
		}},
	})
	if err != nil {
		return nil, err
	}

	matches := []*searchmsg.Match{}
	for _, space := range listSpacesRes.StorageSpaces {
		pathPrefix := ""
		if space.SpaceType == "grant" {
			gpRes, err := p.gwClient.GetPath(ctx, &provider.GetPathRequest{
				ResourceId: space.Root,
			})
			if err != nil {
				return nil, err
			}
			if gpRes.Status.Code != rpcv1beta1.Code_CODE_OK {
				return nil, errtypes.NewErrtypeFromStatus(gpRes.Status)
			}
			pathPrefix = utils.MakeRelativePath(gpRes.Path)
		}

		res, err := p.indexClient.Search(ctx, &searchsvc.SearchIndexRequest{
			Query: req.Query,
			Ref: &searchmsg.Reference{
				ResourceId: &searchmsg.ResourceID{
					StorageId: space.Root.StorageId,
					OpaqueId:  space.Root.OpaqueId,
				},
				Path: pathPrefix,
			},
		})
		if err != nil {
			return nil, err
		}

		for _, match := range res.Matches {
			if pathPrefix != "" {
				match.Entity.Ref.Path = utils.MakeRelativePath(strings.TrimPrefix(match.Entity.Ref.Path, pathPrefix))
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
