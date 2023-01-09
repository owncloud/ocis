package identity

import (
	"context"
	"net/url"

	libregraph "github.com/owncloud/libre-graph-api-go"
)

// GetEducationClasses implements the EducationBackend interface for the LDAP backend.
func (i *LDAP) GetEducationClasses(ctx context.Context, queryParam url.Values) ([]*libregraph.EducationClass, error) {
	return nil, errNotImplemented
}

// CreateEducationClass implements the EducationBackend interface for the LDAP backend.
func (i *LDAP) CreateEducationClass(ctx context.Context, class libregraph.EducationClass) (*libregraph.EducationClass, error) {
	return nil, errNotImplemented
}

// GetEducationClass implements the EducationBackend interface for the LDAP backend.
func (i *LDAP) GetEducationClass(ctx context.Context, namedOrID string, queryParam url.Values) (*libregraph.EducationClass, error) {
	return nil, errNotImplemented
}

// DeleteEducationClass implements the EducationBackend interface for the LDAP backend.
func (i *LDAP) DeleteEducationClass(ctx context.Context, nameOrID string) error {
	return errNotImplemented
}

// GetEducationClassMembers implements the EducationBackend interface for the LDAP backend.
func (i *LDAP) GetEducationClassMembers(ctx context.Context, nameOrID string) ([]*libregraph.EducationUser, error) {
	return nil, errNotImplemented
}
