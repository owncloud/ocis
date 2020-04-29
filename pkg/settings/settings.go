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

type Manager interface {
	BundleManager
	ValueManager
}

// BundleManager
type BundleManager interface {
	ReadBundle(identifier *proto.Identifier) (*proto.SettingsBundle, error)
	WriteBundle(bundle *proto.SettingsBundle) (*proto.SettingsBundle, error)
	ListBundles(identifier *proto.Identifier) ([]*proto.SettingsBundle, error)
}

// ValueManager
type ValueManager interface {
	ReadValue(identifier *proto.Identifier) (*proto.SettingsValue, error)
	WriteValue(value *proto.SettingsValue) (*proto.SettingsValue, error)
	ListValues(identifier *proto.Identifier) ([]*proto.SettingsValue, error)
}
