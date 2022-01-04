package identity

import (
	"context"
	"net/url"

	cs3group "github.com/cs3org/go-cs3apis/cs3/identity/group/v1beta1"
	cs3user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	cs3rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	"github.com/cs3org/reva/pkg/rgrpc/todo/pool"
	libregraph "github.com/owncloud/libre-graph-api-go"

	"github.com/owncloud/ocis/graph/pkg/config"
	"github.com/owncloud/ocis/graph/pkg/service/v0/errorcode"
	"github.com/owncloud/ocis/ocis-pkg/log"
)

type CS3 struct {
	Config *config.Reva
	Logger *log.Logger
}

func (i *CS3) GetUser(ctx context.Context, userID string) (*libregraph.User, error) {
	client, err := pool.GetGatewayServiceClient(i.Config.Address)
	if err != nil {
		i.Logger.Error().Err(err).Msg("could not get client")
		return nil, errorcode.New(errorcode.ServiceNotAvailable, err.Error())
	}

	res, err := client.GetUserByClaim(ctx, &cs3user.GetUserByClaimRequest{
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

func (i *CS3) GetUsers(ctx context.Context, queryParam url.Values) ([]*libregraph.User, error) {
	client, err := pool.GetGatewayServiceClient(i.Config.Address)
	if err != nil {
		i.Logger.Error().Err(err).Msg("could not get client")
		return nil, errorcode.New(errorcode.ServiceNotAvailable, err.Error())
	}

	search := queryParam.Get("search")
	if search == "" {
		search = queryParam.Get("$search")
	}

	res, err := client.FindUsers(ctx, &cs3user.FindUsersRequest{
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

	users := make([]*libregraph.User, 0, len(res.Users))

	for _, user := range res.Users {
		users = append(users, CreateUserModelFromCS3(user))
	}

	return users, nil
}

func (i *CS3) GetGroups(ctx context.Context, queryParam url.Values) ([]*libregraph.Group, error) {
	client, err := pool.GetGatewayServiceClient(i.Config.Address)
	if err != nil {
		i.Logger.Error().Err(err).Msg("could not get client")
		return nil, errorcode.New(errorcode.ServiceNotAvailable, err.Error())
	}

	search := queryParam.Get("search")
	if search == "" {
		search = queryParam.Get("$search")
	}

	res, err := client.FindGroups(ctx, &cs3group.FindGroupsRequest{
		// FIXME presence match is currently not implemented, an empty search currently leads to
		// Unwilling To Perform": Search Error: error parsing filter: (&(objectclass=posixAccount)(|(cn=*)(displayname=*)(mail=*))), error: Present filter match for cn not implemented
		Filter: search,
	})

	switch {
	case err != nil:
		i.Logger.Error().Err(err).Str("search", search).Msg("error sending find groups grpc request")
		return nil, errorcode.New(errorcode.ServiceNotAvailable, err.Error())
	case res.Status.Code != cs3rpc.Code_CODE_OK:
		if res.Status.Code == cs3rpc.Code_CODE_NOT_FOUND {
			return nil, errorcode.New(errorcode.ItemNotFound, res.Status.Message)
		}
		i.Logger.Error().Err(err).Str("search", search).Msg("error sending find groups grpc request")
		return nil, errorcode.New(errorcode.GeneralException, res.Status.Message)
	}

	groups := make([]*libregraph.Group, 0, len(res.Groups))

	for _, group := range res.Groups {
		groups = append(groups, createGroupModelFromCS3(group))
	}

	return groups, nil
}

func (i *CS3) GetGroup(ctx context.Context, groupID string) (*libregraph.Group, error) {
	client, err := pool.GetGatewayServiceClient(i.Config.Address)
	if err != nil {
		i.Logger.Error().Err(err).Msg("could not get client")
		return nil, errorcode.New(errorcode.ServiceNotAvailable, err.Error())
	}

	res, err := client.GetGroupByClaim(ctx, &cs3group.GetGroupByClaimRequest{
		Claim: "groupid", // FIXME add consts to reva
		Value: groupID,
	})

	switch {
	case err != nil:
		i.Logger.Error().Err(err).Str("groupid", groupID).Msg("error sending get group by claim id grpc request")
		return nil, errorcode.New(errorcode.ServiceNotAvailable, err.Error())
	case res.Status.Code != cs3rpc.Code_CODE_OK:
		if res.Status.Code == cs3rpc.Code_CODE_NOT_FOUND {
			return nil, errorcode.New(errorcode.ItemNotFound, res.Status.Message)
		}
		i.Logger.Error().Err(err).Str("groupid", groupID).Msg("error sending get group by claim id grpc request")
		return nil, errorcode.New(errorcode.GeneralException, res.Status.Message)
	}

	return createGroupModelFromCS3(res.Group), nil
}

func createGroupModelFromCS3(g *cs3group.Group) *libregraph.Group {
	if g.Id == nil {
		g.Id = &cs3group.GroupId{}
	}
	return &libregraph.Group{
		Id:                       &g.Id.OpaqueId,
		OnPremisesDomainName:     &g.Id.Idp,
		OnPremisesSamAccountName: &g.GroupName,
		DisplayName:              &g.DisplayName,
		Mail:                     &g.Mail,
		// TODO when to fetch and expand memberof, usernames or ids?
	}
}
