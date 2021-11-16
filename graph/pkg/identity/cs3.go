package identity

import (
	"context"
	"net/url"

	cs3 "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	cs3rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	"github.com/cs3org/reva/pkg/rgrpc/todo/pool"
	msgraph "github.com/yaegashi/msgraph.go/beta"

	"github.com/owncloud/ocis/graph/pkg/config"
	"github.com/owncloud/ocis/graph/pkg/service/v0/errorcode"
	"github.com/owncloud/ocis/ocis-pkg/log"
)

type CS3 struct {
	Config *config.Reva
	Logger *log.Logger
}

func (i *CS3) GetUser(ctx context.Context, userID string) (*msgraph.User, error) {
	client, err := pool.GetGatewayServiceClient(i.Config.Address)
	if err != nil {
		i.Logger.Error().Err(err).Msg("could not get client")
		return nil, errorcode.New(errorcode.ServiceNotAvailable, err.Error())
	}

	res, err := client.GetUserByClaim(ctx, &cs3.GetUserByClaimRequest{
		Claim: "userid", // FIXME add consts to reva
		Value: userID,
	})

	switch {
	case err != nil:
		i.Logger.Error().Err(err).Str("userid", userID).Msg("error sending get user by claim id grpc request")
		return nil, errorcode.New(errorcode.ServiceNotAvailable, err.Error())
	case res.Status.Code != cs3rpc.Code_CODE_OK:
		if res.Status.Code == cs3rpc.Code_CODE_NOT_FOUND {
			return nil, errorcode.New(errorcode.ItemNotFound, res.Status.Message)
		}
		i.Logger.Error().Err(err).Str("userid", userID).Msg("error sending get user by claim id grpc request")
		return nil, errorcode.New(errorcode.GeneralException, res.Status.Message)
	}
	return CreateUserModelFromCS3(res.User), nil
}

func (i *CS3) GetUsers(ctx context.Context, queryParam url.Values) ([]*msgraph.User, error) {
	client, err := pool.GetGatewayServiceClient(i.Config.Address)
	if err != nil {
		i.Logger.Error().Err(err).Msg("could not get client")
		return nil, errorcode.New(errorcode.ServiceNotAvailable, err.Error())
	}

	search := queryParam.Get("search")
	if search == "" {
		search = queryParam.Get("$search")
	}

	res, err := client.FindUsers(ctx, &cs3.FindUsersRequest{
		// FIXME presence match is currently not implemented, an empty search currently leads to
		// Unwilling To Perform": Search Error: error parsing filter: (&(objectclass=posixAccount)(|(cn=*)(displayname=*)(mail=*))), error: Present filter match for cn not implemented
		Filter: search,
	})
	switch {
	case err != nil:
		i.Logger.Error().Err(err).Str("search", search).Msg("error sending find users grpc request")
		return nil, errorcode.New(errorcode.ServiceNotAvailable, err.Error())
	case res.Status.Code != cs3rpc.Code_CODE_OK:
		if res.Status.Code == cs3rpc.Code_CODE_NOT_FOUND {
			return nil, errorcode.New(errorcode.ItemNotFound, res.Status.Message)
		}
		i.Logger.Error().Err(err).Str("search", search).Msg("error sending find users grpc request")
		return nil, errorcode.New(errorcode.GeneralException, res.Status.Message)
	}

	users := make([]*msgraph.User, 0, len(res.Users))

	for _, user := range res.Users {
		users = append(users, CreateUserModelFromCS3(user))
	}

	return users, nil
}
