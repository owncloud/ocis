// Package store implements the go-micro store interface
package store

import (
	"errors"

	settingsmsg "github.com/owncloud/ocis/protogen/gen/ocis/messages/settings/v0"
)

// ListRoleAssignments loads and returns all role assignments matching the given assignment identifier.
func (s Store) ListRoleAssignments(accountUUID string) ([]*settingsmsg.UserRoleAssignment, error) {
	return nil, errors.New("not implemented")
}

// WriteRoleAssignment appends the given role assignment to the existing assignments of the respective account.
func (s Store) WriteRoleAssignment(accountUUID, roleID string) (*settingsmsg.UserRoleAssignment, error) {
	return nil, errors.New("not implemented")
}

// RemoveRoleAssignment deletes the given role assignment from the existing assignments of the respective account.
func (s Store) RemoveRoleAssignment(assignmentID string) error {
	return errors.New("not implemented")
}
