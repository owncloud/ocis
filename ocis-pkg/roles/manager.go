package roles

import (
	"context"

	"github.com/owncloud/ocis/ocis-pkg/log"
	settings "github.com/owncloud/ocis/settings/pkg/proto/v0"
)

// Manager manages a cache of roles by fetching unknown roles from the settings.RoleService.
type Manager struct {
	logger      log.Logger
	cache       cache
	roleService settings.RoleService
}

// NewManager returns a new instance of Manager.
func NewManager(o ...Option) Manager {
	opts := newOptions(o...)

	return Manager{
		cache:       newCache(opts.size, opts.ttl),
		roleService: opts.roleService,
	}
}

// List returns all roles that match the given roleIDs.
func (m *Manager) List(ctx context.Context, roleIDs []string) []*settings.Bundle {
	// get from cache
	result := make([]*settings.Bundle, 0)
	lookup := make([]string, 0)
	for _, roleID := range roleIDs {
		if hit := m.cache.get(roleID); hit == nil {
			lookup = append(lookup, roleID)
		} else {
			result = append(result, hit)
		}
	}

	// if there are roles missing, fetch them from the RoleService
	if len(lookup) > 0 {
		request := &settings.ListBundlesRequest{
			BundleIds: lookup,
		}
		res, err := m.roleService.ListRoles(ctx, request)
		if err != nil {
			m.logger.Debug().Err(err).Msg("failed to fetch roles by roleIDs")
		}
		for _, role := range res.Bundles {
			m.cache.set(role.Id, role)
			result = append(result, role)
		}
	}

	return result
}

// FindPermissionByID searches for a permission-setting by the permissionID, but limited to the given roleIDs
func (m *Manager) FindPermissionByID(ctx context.Context, roleIDs []string, permissionID string) *settings.Setting {
	for _, role := range m.List(ctx, roleIDs) {
		for _, setting := range role.Settings {
			if setting.Id == permissionID {
				return setting
			}
		}
	}
	return nil
}

