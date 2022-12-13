package identity

import (
	"context"
	"net/url"

	libregraph "github.com/owncloud/libre-graph-api-go"
)

// NOOP is a dummy EducationBackend, doing nothing
type NOOP struct{}

// CreateSchool creates the supplied school in the identity backend.
func (i *NOOP) CreateSchool(ctx context.Context, school libregraph.EducationSchool) (*libregraph.EducationSchool, error) {
	return nil, errNotImplemented
}

// DeleteSchool deletes a given school, identified by id
func (i *NOOP) DeleteSchool(ctx context.Context, id string) error {
	return errNotImplemented
}

// GetSchool implements the EducationBackend interface for the NOOP backend.
func (i *NOOP) GetSchool(ctx context.Context, nameOrID string, queryParam url.Values) (*libregraph.EducationSchool, error) {
	return nil, errNotImplemented
}

// GetSchools implements the EducationBackend interface for the NOOP backend.
func (i *NOOP) GetSchools(ctx context.Context, queryParam url.Values) ([]*libregraph.EducationSchool, error) {
	return nil, errNotImplemented
}

// GetSchoolMembers implements the EducationBackend interface for the NOOP backend.
func (i *NOOP) GetSchoolMembers(ctx context.Context, id string) ([]*libregraph.EducationUser, error) {
	return nil, errNotImplemented
}

// AddMembersToSchool adds new members (reference by a slice of IDs) to supplied school in the identity backend.
func (i *NOOP) AddMembersToSchool(ctx context.Context, schoolID string, memberID []string) error {
	return errNotImplemented
}

// RemoveMemberFromSchool removes a single member (by ID) from a school
func (i *NOOP) RemoveMemberFromSchool(ctx context.Context, schoolID string, memberID string) error {
	return errNotImplemented
}

// CreateEducationUser creates a given education user in the identity backend.
func (i *NOOP) CreateEducationUser(ctx context.Context, user libregraph.EducationUser) (*libregraph.EducationUser, error) {
	return nil, errNotImplemented
}

// DeleteEducationUser deletes a given educationuser, identified by username or id, from the backend
func (i *NOOP) DeleteEducationUser(ctx context.Context, nameOrID string) error {
	return errNotImplemented
}

// UpdateEducationUser applies changes to given education user, identified by username or id
func (i *NOOP) UpdateEducationUser(ctx context.Context, nameOrID string, user libregraph.EducationUser) (*libregraph.EducationUser, error) {
	return nil, errNotImplemented
}

// GetEducationUser implements the EducationBackend interface for the NOOP backend.
func (i *NOOP) GetEducationUser(ctx context.Context, nameOrID string, queryParam url.Values) (*libregraph.EducationUser, error) {
	return nil, errNotImplemented
}

// GetEducationUsers implements the EducationBackend interface for the NOOP backend.
func (i *NOOP) GetEducationUsers(ctx context.Context, queryParam url.Values) ([]*libregraph.EducationUser, error) {
	return nil, errNotImplemented
}
