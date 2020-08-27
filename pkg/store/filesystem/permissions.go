package store

import (
	"github.com/owncloud/ocis-settings/pkg/proto/v0"
	"github.com/owncloud/ocis-settings/pkg/util"
)

// ListPermissionsByResource collects all permissions from the provided roleIDs that match the requested resource
func (s Store) ListPermissionsByResource(resource *proto.Resource, roleIDs []string) ([]*proto.Permission, error) {
	records := make([]*proto.Permission, 0)
	for _, roleID := range roleIDs {
		role, err := s.ReadBundle(roleID)
		if err != nil {
			s.Logger.Debug().Str("roleID", roleID).Msg("role not found, skipping")
			continue
		}
		records = append(records, extractPermissionsByResource(resource, role)...)
	}
	return records, nil
}

// ReadPermissionByID finds the permission in the roles, specified by the provided roleIDs
func (s Store) ReadPermissionByID(permissionID string, roleIDs []string) (*proto.Permission, error) {
	for _, roleID := range roleIDs {
		role, err := s.ReadBundle(roleID)
		if err != nil {
			s.Logger.Debug().Str("roleID", roleID).Msg("role not found, skipping")
			continue
		}
		for _, permission := range role.Settings {
			if permission.Id == permissionID {
				if value, ok := permission.Value.(*proto.Setting_PermissionValue); ok {
					return value.PermissionValue, nil
				}
			}
		}
	}
	return nil, nil
}

// extractPermissionsByResource collects all permissions from the provided role that match the requested resource
func extractPermissionsByResource(resource *proto.Resource, role *proto.Bundle) []*proto.Permission {
	permissions := make([]*proto.Permission, 0)
	for _, setting := range role.Settings {
		if value, ok := setting.Value.(*proto.Setting_PermissionValue); ok {
			if util.IsResourceMatched(setting.Resource, resource) {
				permissions = append(permissions, value.PermissionValue)
			}
		}
	}
	return permissions
}
