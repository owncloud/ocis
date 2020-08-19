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
}

// BundleManager is a bundle service interface for abstraction of storage implementations
type BundleManager interface {
	ListBundles(bundleType proto.Bundle_Type) ([]*proto.Bundle, error)
	ReadBundle(bundleID string) (*proto.Bundle, error)
	WriteBundle(bundle *proto.Bundle) (*proto.Bundle, error)
	ReadSetting(settingID string) (*proto.Setting, error)
}

// ValueManager is a value service interface for abstraction of storage implementations
type ValueManager interface {
	ListValues(bundleID, accountUUID string) ([]*proto.Value, error)
	ReadValue(valueID string) (*proto.Value, error)
	ReadValueByUniqueIdentifiers(accountUUID, settingID string) (*proto.Value, error)
	WriteValue(value *proto.Value) (*proto.Value, error)
}
