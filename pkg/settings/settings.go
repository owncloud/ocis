package settings

import (
	"github.com/owncloud/ocis-settings/pkg/config"
	"github.com/owncloud/ocis-settings/pkg/proto/v0"
)

var (
	// Registry uses the strategy pattern as a registry
	Registry = map[string]RegisterFunc{}
)

// RegisterFunc stores store constructors
type RegisterFunc func(*config.Config) Manager

// Manager combines service interfaces for abstraction of storage implementations
type Manager interface {
	BundleManager
	ValueManager
	RoleAssignmentManager
	PermissionManager
}

// BundleManager is a bundle service interface for abstraction of storage implementations
type BundleManager interface {
	ListBundles(bundleType proto.Bundle_Type, bundleIDs []string) ([]*proto.Bundle, error)
	ReadBundle(bundleID string) (*proto.Bundle, error)
	WriteBundle(bundle *proto.Bundle) (*proto.Bundle, error)
	ReadSetting(settingID string) (*proto.Setting, error)
	AddSettingToBundle(bundleID string, setting *proto.Setting) (*proto.Setting, error)
	RemoveSettingFromBundle(bundleID, settingID string) error
}

// ValueManager is a value service interface for abstraction of storage implementations
type ValueManager interface {
	ListValues(bundleID, accountUUID string) ([]*proto.Value, error)
	ReadValue(valueID string) (*proto.Value, error)
	ReadValueByUniqueIdentifiers(accountUUID, settingID string) (*proto.Value, error)
	WriteValue(value *proto.Value) (*proto.Value, error)
}

// RoleAssignmentManager is a role assignment service interface for abstraction of storage implementations
type RoleAssignmentManager interface {
	ListRoleAssignments(accountUUID string) ([]*proto.UserRoleAssignment, error)
	WriteRoleAssignment(accountUUID, roleID string) (*proto.UserRoleAssignment, error)
	RemoveRoleAssignment(assignmentID string) error
}

// PermissionManager is a permissions service interface for abstraction of storage implementations
type PermissionManager interface {
	ListPermissionsByResource(resource *proto.Resource, roleIDs []string) ([]*proto.Permission, error)
	ReadPermissionByID(permissionID string, roleIDs []string) (*proto.Permission, error)
}
