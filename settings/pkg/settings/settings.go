package settings

import (
	"errors"

	settingsmsg "github.com/owncloud/ocis/protogen/gen/ocis/messages/settings/v0"
	"github.com/owncloud/ocis/settings/pkg/config"
)

var (
	// Registry uses the strategy pattern as a registry
	Registry = map[string]RegisterFunc{}

	// ErrPermissionNotFound defines a new error for when a permission was not found
	ErrPermissionNotFound = errors.New("permission not found")
)

// RegisterFunc stores store constructors
type RegisterFunc func(*config.Config, func(Manager)) Manager

// Manager combines service interfaces for abstraction of storage implementations
type Manager interface {
	BundleManager
	ValueManager
	RoleAssignmentManager
	PermissionManager
}

// BundleManager is a bundle service interface for abstraction of storage implementations
type BundleManager interface {
	ListBundles(bundleType settingsmsg.Bundle_Type, bundleIDs []string) ([]*settingsmsg.Bundle, error)
	ReadBundle(bundleID string) (*settingsmsg.Bundle, error)
	WriteBundle(bundle *settingsmsg.Bundle) (*settingsmsg.Bundle, error)
	ReadSetting(settingID string) (*settingsmsg.Setting, error)
	AddSettingToBundle(bundleID string, setting *settingsmsg.Setting) (*settingsmsg.Setting, error)
	RemoveSettingFromBundle(bundleID, settingID string) error
}

// ValueManager is a value service interface for abstraction of storage implementations
type ValueManager interface {
	ListValues(bundleID, accountUUID string) ([]*settingsmsg.Value, error)
	ReadValue(valueID string) (*settingsmsg.Value, error)
	ReadValueByUniqueIdentifiers(accountUUID, settingID string) (*settingsmsg.Value, error)
	WriteValue(value *settingsmsg.Value) (*settingsmsg.Value, error)
}

// RoleAssignmentManager is a role assignment service interface for abstraction of storage implementations
type RoleAssignmentManager interface {
	ListRoleAssignments(accountUUID string) ([]*settingsmsg.UserRoleAssignment, error)
	WriteRoleAssignment(accountUUID, roleID string) (*settingsmsg.UserRoleAssignment, error)
	RemoveRoleAssignment(assignmentID string) error
}

// PermissionManager is a permissions service interface for abstraction of storage implementations
type PermissionManager interface {
	ListPermissionsByResource(resource *settingsmsg.Resource, roleIDs []string) ([]*settingsmsg.Permission, error)
	ReadPermissionByID(permissionID string, roleIDs []string) (*settingsmsg.Permission, error)
	ReadPermissionByName(name string, roleIDs []string) (*settingsmsg.Permission, error)
}
