package store

import (
	"errors"

	settingsmsg "github.com/owncloud/ocis/protogen/gen/ocis/messages/settings/v0"
)

// ListPermissionsByResource collects all permissions from the provided roleIDs that match the requested resource
func (s Store) ListPermissionsByResource(resource *settingsmsg.Resource, roleIDs []string) ([]*settingsmsg.Permission, error) {
	return nil, errors.New("not implemented")
}

// ReadPermissionByID finds the permission in the roles, specified by the provided roleIDs
func (s Store) ReadPermissionByID(permissionID string, roleIDs []string) (*settingsmsg.Permission, error) {
	return nil, errors.New("not implemented")
}

// ReadPermissionByName finds the permission in the roles, specified by the provided roleIDs
func (s Store) ReadPermissionByName(name string, roleIDs []string) (*settingsmsg.Permission, error) {
	return nil, errors.New("not implemented")
}
