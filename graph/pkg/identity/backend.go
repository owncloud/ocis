package identity

import (
	"context"
	"net/url"

	cs3 "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	libregraph "github.com/owncloud/libre-graph-api-go"
)

type Backend interface {
	// CreateUser creates a given user in the identity backend.
	CreateUser(ctx context.Context, user libregraph.User) (*libregraph.User, error)

	// DeleteUser deletes a given user, identified by username or id, from the backend
	DeleteUser(ctx context.Context, nameOrID string) error

	// UpdateUser applies changes to given user, identified by username or id
	UpdateUser(ctx context.Context, nameOrID string, user libregraph.User) (*libregraph.User, error)

	GetUser(ctx context.Context, nameOrID string) (*libregraph.User, error)
	GetUsers(ctx context.Context, queryParam url.Values) ([]*libregraph.User, error)

	// CreateGroup creates the supplied group in the identity backend.
	CreateGroup(ctx context.Context, group libregraph.Group) (*libregraph.Group, error)
	// DeleteGroup deletes a given group, identified by id
	DeleteGroup(ctx context.Context, id string) error
	GetGroup(ctx context.Context, nameOrID string) (*libregraph.Group, error)
	GetGroups(ctx context.Context, queryParam url.Values) ([]*libregraph.Group, error)
	GetGroupMembers(ctx context.Context, id string) ([]*libregraph.User, error)
	// AddMemberToGroup adds a new member (reference by ID) to supplied group in the identity backend.
	AddMemberToGroup(ctx context.Context, groupID string, memberID string) error
	// RemoveMemberFromGroup removes a single member (by ID) from a group
	RemoveMemberFromGroup(ctx context.Context, groupID string, memberID string) error
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
