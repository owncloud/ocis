package search

import (
	"context"
	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	ctxpkg "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/cs3org/reva/v2/pkg/errtypes"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	searchmsg "github.com/owncloud/ocis/v2/protogen/gen/ocis/messages/search/v0"
	"github.com/owncloud/ocis/v2/services/search/pkg/engine"
	"google.golang.org/grpc/metadata"
)

type MatchArray []*searchmsg.Match

func (s MatchArray) Len() int {
	return len(s)
}
func (s MatchArray) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s MatchArray) Less(i, j int) bool {
	return s[i].Score > s[j].Score
}

func logDocCount(engine engine.Engine, logger log.Logger) {
	c, err := engine.DocCount()
	if err != nil {
		logger.Error().Err(err).Msg("error getting document count from the index")
	}
	logger.Debug().Interface("count", c).Msg("new document count")
}

func getAuthContext(owner *user.User, gw gateway.GatewayAPIClient, secret string, logger log.Logger) (context.Context, error) {
	ownerCtx := ctxpkg.ContextSetUser(context.Background(), owner)
	authRes, err := gw.Authenticate(ownerCtx, &gateway.AuthenticateRequest{
		Type:         "machine",
		ClientId:     "userid:" + owner.GetId().GetOpaqueId(),
		ClientSecret: secret,
	})

	if err == nil && authRes.GetStatus().GetCode() != rpc.Code_CODE_OK {
		err = errtypes.NewErrtypeFromStatus(authRes.Status)
	}

	if err != nil {
		logger.Error().Err(err).Interface("owner", owner).Interface("authRes", authRes).Msg("error using machine auth")
		return nil, err
	}

	return metadata.AppendToOutgoingContext(ownerCtx, ctxpkg.TokenHeader, authRes.Token), nil
}

func statResource(ref *provider.Reference, owner *user.User, gw gateway.GatewayAPIClient, secret string, logger log.Logger) (*provider.StatResponse, error) {
	ownerCtx, err := getAuthContext(owner, gw, secret, logger)
	if err != nil {
		return nil, err
	}

	return gw.Stat(ownerCtx, &provider.StatRequest{Ref: ref})
}

func getPath(id *provider.ResourceId, owner *user.User, gw gateway.GatewayAPIClient, secret string, logger log.Logger) (*provider.GetPathResponse, error) {
	ownerCtx, err := getAuthContext(owner, gw, secret, logger)
	if err != nil {
		return nil, err
	}

	return gw.GetPath(ownerCtx, &provider.GetPathRequest{ResourceId: id})
}
