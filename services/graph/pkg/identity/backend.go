package identity

import (
	"context"
	"net/url"

	cs3 "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	libregraph "github.com/owncloud/libre-graph-api-go"
)

// Backend defines the Interface for an IdentityBackend implementation
type Backend interface {
	// CreateUser creates a given user in the identity backend.
	CreateUser(ctx context.Context, user libregraph.User) (*libregraph.User, error)
	// DeleteUser deletes a given user, identified by username or id, from the backend
	DeleteUser(ctx context.Context, nameOrID string) error
	// UpdateUser applies changes to given user, identified by username or id
	UpdateUser(ctx context.Context, nameOrID string, user libregraph.User) (*libregraph.User, error)
	GetUser(ctx context.Context, nameOrID string, queryParam url.Values) (*libregraph.User, error)
	GetUsers(ctx context.Context, queryParam url.Values) ([]*libregraph.User, error)

	// CreateGroup creates the supplied group in the identity backend.
	CreateGroup(ctx context.Context, group libregraph.Group) (*libregraph.Group, error)
	// DeleteGroup deletes a given group, identified by id
	DeleteGroup(ctx context.Context, id string) error
	GetGroup(ctx context.Context, nameOrID string, queryParam url.Values) (*libregraph.Group, error)
	GetGroups(ctx context.Context, queryParam url.Values) ([]*libregraph.Group, error)
	GetGroupMembers(ctx context.Context, id string) ([]*libregraph.User, error)
	// AddMembersToGroup adds new members (reference by a slice of IDs) to supplied group in the identity backend.
	AddMembersToGroup(ctx context.Context, groupID string, memberID []string) error
	// RemoveMemberFromGroup removes a single member (by ID) from a group
	RemoveMemberFromGroup(ctx context.Context, groupID string, memberID string) error
}

// EducationBackend defines the Interface for an EducationBackend implementation
type EducationBackend interface {
	// CreateSchool creates the supplied school in the identity backend.
	CreateSchool(ctx context.Context, group libregraph.EducationSchool) (*libregraph.EducationSchool, error)
	// DeleteSchool deletes a given school, identified by id
	DeleteSchool(ctx context.Context, id string) error
	// GetSchool reads a given school by id
	GetSchool(ctx context.Context, nameOrID string, queryParam url.Values) (*libregraph.EducationSchool, error)
	// GetSchools lists all	schools
	GetSchools(ctx context.Context, queryParam url.Values) ([]*libregraph.EducationSchool, error)
	GetSchoolMembers(ctx context.Context, id string) ([]*libregraph.User, error)
	// AddMembersToSchool adds new members (reference by a slice of IDs) to supplied school in the identity backend.
	AddMembersToSchool(ctx context.Context, schoolID string, memberID []string) error
	// RemoveMemberFromSchool removes a single member (by ID) from a school
	RemoveMemberFromSchool(ctx context.Context, schoolID string, memberID string) error
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
