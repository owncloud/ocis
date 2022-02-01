package store

import (
	settingsmsg "github.com/owncloud/ocis/protogen/gen/ocis/messages/settings/v0"
	"github.com/owncloud/ocis/settings/pkg/settings"
	"github.com/owncloud/ocis/settings/pkg/util"
)

// ListPermissionsByResource collects all permissions from the provided roleIDs that match the requested resource
func (s Store) ListPermissionsByResource(resource *settingsmsg.Resource, roleIDs []string) ([]*settingsmsg.Permission, error) {
	records := make([]*settingsmsg.Permission, 0)
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
func (s Store) ReadPermissionByID(permissionID string, roleIDs []string) (*settingsmsg.Permission, error) {
	for _, roleID := range roleIDs {
		role, err := s.ReadBundle(roleID)
		if err != nil {
			s.Logger.Debug().Str("roleID", roleID).Msg("role not found, skipping")
			continue
		}
		for _, permission := range role.Settings {
			if permission.Id == permissionID {
				if value, ok := permission.Value.(*settingsmsg.Setting_PermissionValue); ok {
					return value.PermissionValue, nil
				}
			}
		}
	}
	return nil, nil
}

// ReadPermissionByName finds the permission in the roles, specified by the provided roleIDs
func (s Store) ReadPermissionByName(name string, roleIDs []string) (*settingsmsg.Permission, error) {
	for _, roleID := range roleIDs {
		role, err := s.ReadBundle(roleID)
		if err != nil {
			s.Logger.Debug().Str("roleID", roleID).Msg("role not found, skipping")
			continue
		}
		for _, permission := range role.Settings {
			if permission.Name == name {
				if value, ok := permission.Value.(*settingsmsg.Setting_PermissionValue); ok {
					return value.PermissionValue, nil
				}
			}
		}
	}
	return nil, settings.ErrPermissionNotFound
}

// extractPermissionsByResource collects all permissions from the provided role that match the requested resource
func extractPermissionsByResource(resource *settingsmsg.Resource, role *settingsmsg.Bundle) []*settingsmsg.Permission {
	permissions := make([]*settingsmsg.Permission, 0)
	for _, setting := range role.Settings {
		if value, ok := setting.Value.(*settingsmsg.Setting_PermissionValue); ok {
			if util.IsResourceMatched(setting.Resource, resource) {
				permissions = append(permissions, value.PermissionValue)
			}
		}
	}
	return permissions
}
