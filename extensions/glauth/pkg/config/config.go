package config

import (
	"context"

	"github.com/owncloud/ocis/ocis-pkg/shared"
)

// Config combines all available configuration parts.
type Config struct {
	*shared.Commons `yaml:"-"`

	Service Service `yaml:"-"`

	Tracing *Tracing `yaml:"tracing"`
	Log     *Log     `yaml:"log"`
	Debug   Debug    `yaml:"debug"`

	Ldap  Ldap  `yaml:"ldap"`
	Ldaps Ldaps `yaml:"ldaps"`

	Backend  Backend         `yaml:"backend"`
	Fallback FallbackBackend `yaml:"fallback"`

	RoleBundleUUID string `yaml:"role_bundle_uuid" env:"GLAUTH_ROLE_BUNDLE_ID"`

	ConfigFile string `yaml:"-" env:"GLAUTH_CONFIG_FILE" desc:"config file to be used by the glauth extension"`

	Context context.Context `yaml:"-"`
}

// Backend defined the available backend configuration.
type Backend struct {
	Datastore   string   `yaml:"datastore"`
	BaseDN      string   `yaml:"base_dn"`
	Insecure    bool     `yaml:"insecure"`
	NameFormat  string   `yaml:"name_format"`
	GroupFormat string   `yaml:"group_format"`
	Servers     []string `yaml:"servers"`
	SSHKeyAttr  string   `yaml:"ssh_key_attr"`
	UseGraphAPI bool     `yaml:"use_graph_api"`
}

// FallbackBackend defined the available fallback backend configuration.
type FallbackBackend struct {
	Datastore   string   `yaml:"datastore"`
	BaseDN      string   `yaml:"base_dn"`
	Insecure    bool     `yaml:"insecure"`
	NameFormat  string   `yaml:"name_format"`
	GroupFormat string   `yaml:"group_format"`
	Servers     []string `yaml:"servers"`
	SSHKeyAttr  string   `yaml:"ssh_key_attr"`
	UseGraphAPI bool     `yaml:"use_graph_api"`
}
