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

// Manager
type Manager interface {
	Read(extension string, key string) (*proto.SettingsBundle, error)
	Write(bundle *proto.SettingsBundle) (*proto.SettingsBundle, error)
	ListAll() ([]*proto.SettingsBundle, error)
	ListByExtension(extension string) ([]*proto.SettingsBundle, error)
}
