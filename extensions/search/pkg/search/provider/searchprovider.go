package provider

import (
	"context"
	"strings"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	rpcv1beta1 "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	providerv1beta1 "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	typesv1beta1 "github.com/cs3org/go-cs3apis/cs3/types/v1beta1"
	ctxpkg "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/cs3org/reva/v2/pkg/errtypes"
	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/cs3org/reva/v2/pkg/utils"
	"github.com/owncloud/ocis/extensions/search/pkg/search"
	"google.golang.org/grpc/metadata"

	searchmsg "github.com/owncloud/ocis/protogen/gen/ocis/messages/search/v0"
	searchsvc "github.com/owncloud/ocis/protogen/gen/ocis/services/search/v0"
)

type Provider struct {
	gwClient          gateway.GatewayAPIClient
	indexClient       search.IndexClient
	machineAuthAPIKey string
}

func New(gwClient gateway.GatewayAPIClient, indexClient search.IndexClient, machineAuthAPIKey string, eventsChan <-chan interface{}) *Provider {
	go func() {
		for {
			ev := <-eventsChan
			var ref *providerv1beta1.Reference
			var owner *user.User
			switch e := ev.(type) {
			case events.FileUploaded:
				ref = e.FileID
				owner = &user.User{
					Id: e.Executant,
				}
			default:
				// Not sure what to do here. Skip.
				continue
			}

			// Get auth
			ownerCtx := ctxpkg.ContextSetUser(context.Background(), owner)
			authRes, err := gwClient.Authenticate(ownerCtx, &gateway.AuthenticateRequest{
				Type:         "machine",
				ClientId:     "userid:" + owner.Id.OpaqueId,
				ClientSecret: machineAuthAPIKey,
			})
			if err != nil || authRes.GetStatus().GetCode() != rpc.Code_CODE_OK {
				// TODO: log error
			}
			ownerCtx = metadata.AppendToOutgoingContext(ownerCtx, ctxpkg.TokenHeader, authRes.Token)

			// Stat changed resource resource
			statRes, err := gwClient.Stat(ownerCtx, &providerv1beta1.StatRequest{Ref: ref})
			if err != nil || statRes.Status.Code != rpc.Code_CODE_OK {
				// TODO: log error
			}

			indexClient.Add(ref, statRes.Info)
		}
	}()

	return &Provider{
		gwClient:          gwClient,
		indexClient:       indexClient,
		machineAuthAPIKey: machineAuthAPIKey,
	}
}

func (p *Provider) Search(ctx context.Context, req *searchsvc.SearchRequest) (*searchsvc.SearchResponse, error) {
	if req.Query == "" {
		return nil, errtypes.PreconditionFailed("empty query provided")
	}

	listSpacesRes, err := p.gwClient.ListStorageSpaces(ctx, &providerv1beta1.ListStorageSpacesRequest{
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
			gpRes, err := p.gwClient.GetPath(ctx, &providerv1beta1.GetPathRequest{
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
