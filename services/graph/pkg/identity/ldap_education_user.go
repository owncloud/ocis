package identity

import (
	"context"
	"net/url"

	libregraph "github.com/owncloud/libre-graph-api-go"
)

// CreateEducationUser creates a given education user in the identity backend.
func (i *LDAP) CreateEducationUser(ctx context.Context, user libregraph.EducationUser) (*libregraph.EducationUser, error) {
	return nil, errNotImplemented
}

// DeleteEducationUser deletes a given educationuser, identified by username or id, from the backend
func (i *LDAP) DeleteEducationUser(ctx context.Context, nameOrID string) error {
	return errNotImplemented
}

// UpdateEducationUser applies changes to given education user, identified by username or id
func (i *LDAP) UpdateEducationUser(ctx context.Context, nameOrID string, user libregraph.EducationUser) (*libregraph.EducationUser, error) {
	return nil, errNotImplemented
}

// GetEducationUser implements the EducationBackend interface for the LDAP backend.
func (i *LDAP) GetEducationUser(ctx context.Context, nameOrID string, queryParam url.Values) (*libregraph.EducationUser, error) {
	return nil, errNotImplemented
}

// GetEducationUsers implements the EducationBackend interface for the LDAP backend.
func (i *LDAP) GetEducationUsers(ctx context.Context, queryParam url.Values) ([]*libregraph.EducationUser, error) {
	return nil, errNotImplemented
}
