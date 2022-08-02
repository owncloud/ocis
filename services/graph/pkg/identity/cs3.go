package identity

import (
	"context"
	"net/url"

	cs3group "github.com/cs3org/go-cs3apis/cs3/identity/group/v1beta1"
	cs3user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	cs3rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	"github.com/cs3org/reva/v2/pkg/rgrpc/todo/pool"
	libregraph "github.com/owncloud/libre-graph-api-go"

	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/graph/pkg/config"
	"github.com/owncloud/ocis/v2/services/graph/pkg/service/v0/errorcode"
)

var (
	errNotImplemented = errorcode.New(errorcode.NotSupported, "not implemented")
)

type CS3 struct {
	Config *config.Reva
	Logger *log.Logger
}

// CreateUser implements the Backend Interface. It's currently not supported for the CS3 backend
func (i *CS3) CreateUser(ctx context.Context, user libregraph.User) (*libregraph.User, error) {
	return nil, errNotImplemented
}

// DeleteUser implements the Backend Interface. It's currently not supported for the CS3 backend
func (i *CS3) DeleteUser(ctx context.Context, nameOrID string) error {
	return errNotImplemented
}

// UpdateUser implements the Backend Interface. It's currently not suported for the CS3 backend
func (i *CS3) UpdateUser(ctx context.Context, nameOrID string, user libregraph.User) (*libregraph.User, error) {
	return nil, errNotImplemented
}

func (i *CS3) GetUser(ctx context.Context, userID string, queryParam url.Values) (*libregraph.User, error) {
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

// CreateGroup implements the Backend Interface. It's currently not supported for the CS3 backend
func (i *CS3) CreateGroup(ctx context.Context, group libregraph.Group) (*libregraph.Group, error) {
	return nil, errorcode.New(errorcode.NotSupported, "not implemented")
}

func (i *CS3) GetGroup(ctx context.Context, groupID string, queryParam url.Values) (*libregraph.Group, error) {
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

// DeleteGroup implements the Backend Interface. It's currently not supported for the CS3 backend
func (i *CS3) DeleteGroup(ctx context.Context, id string) error {
	return errorcode.New(errorcode.NotSupported, "not implemented")
}

// GetGroupMembers implements the Backend Interface. It's currently not supported for the CS3 backend
func (i *CS3) GetGroupMembers(ctx context.Context, groupID string) ([]*libregraph.User, error) {
	return nil, errorcode.New(errorcode.NotSupported, "not implemented")
}

// AddMembersToGroup implements the Backend Interface. It's currently not supported for the CS3 backend
func (i *CS3) AddMembersToGroup(ctx context.Context, groupID string, memberID []string) error {
	return errorcode.New(errorcode.NotSupported, "not implemented")
}

// RemoveMemberFromGroup implements the Backend Interface. It's currently not supported for the CS3 backend
func (i *CS3) RemoveMemberFromGroup(ctx context.Context, groupID string, memberID string) error {
	return errorcode.New(errorcode.NotSupported, "not implemented")
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
		// TODO when to fetch and expand memberof, usernames or ids?
	}
}
