package identity

import (
	"context"
	"net/url"

	cs3 "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	libregraph "github.com/owncloud/libre-graph-api-go"
)

type Backend interface {
	GetUser(ctx context.Context, nameOrId string) (*libregraph.User, error)
	GetUsers(ctx context.Context, queryParam url.Values) ([]*libregraph.User, error)

	GetGroup(ctx context.Context, nameOrId string) (*libregraph.Group, error)
	GetGroups(ctx context.Context, queryParam url.Values) ([]*libregraph.Group, error)
}

func CreateUserModelFromCS3(u *cs3.User) *libregraph.User {
	if u.Id == nil {
		u.Id = &cs3.UserId{}
	}
	return &libregraph.User{
		DisplayName:              &u.DisplayName,
		Mail:                     &u.Mail,
		OnPremisesSamAccountName: &u.Username,
		Id:                       &u.Id.OpaqueId,
	}
}
