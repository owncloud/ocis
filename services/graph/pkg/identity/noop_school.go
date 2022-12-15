package identity

import (
	"context"
	"net/url"

	libregraph "github.com/owncloud/libre-graph-api-go"
)

// NOOP is a dummy EducationBackend, doing nothing
type NOOP struct{}

// CreateEducationSchool creates the supplied school in the identity backend.
func (i *NOOP) CreateEducationSchool(ctx context.Context, school libregraph.EducationSchool) (*libregraph.EducationSchool, error) {
	return nil, errNotImplemented
}

// DeleteEducationSchool deletes a given school, identified by id
func (i *NOOP) DeleteEducationSchool(ctx context.Context, id string) error {
	return errNotImplemented
}

// GetEducationSchool implements the EducationBackend interface for the NOOP backend.
func (i *NOOP) GetEducationSchool(ctx context.Context, nameOrID string, queryParam url.Values) (*libregraph.EducationSchool, error) {
	return nil, errNotImplemented
}

// GetEducationSchools implements the EducationBackend interface for the NOOP backend.
func (i *NOOP) GetEducationSchools(ctx context.Context, queryParam url.Values) ([]*libregraph.EducationSchool, error) {
	return nil, errNotImplemented
}

// GetEducationSchoolMembers implements the EducationBackend interface for the NOOP backend.
func (i *NOOP) GetEducationSchoolMembers(ctx context.Context, id string) ([]*libregraph.EducationUser, error) {
	return nil, errNotImplemented
}

// AddMembersToEducationSchool adds new members (reference by a slice of IDs) to supplied school in the identity backend.
func (i *NOOP) AddMembersToEducationSchool(ctx context.Context, schoolID string, memberID []string) error {
	return errNotImplemented
}

// RemoveMemberFromEducationSchool removes a single member (by ID) from a school
func (i *NOOP) RemoveMemberFromEducationSchool(ctx context.Context, schoolID string, memberID string) error {
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
