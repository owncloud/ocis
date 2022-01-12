package config

import (
	"context"

	"github.com/owncloud/ocis/ocis-pkg/shared"
)

// Config combines all available configuration parts.
type Config struct {
	*shared.Commons

	Service Service

	Tracing *Tracing `ocisConfig:"tracing"`
	Log     *Log     `ocisConfig:"log"`
	Debug   Debug    `ocisConfig:"debug"`

	Ldap  Ldap  `ocisConfig:"ldap"`
	Ldaps Ldaps `ocisConfig:"ldaps"`

	Backend  Backend         `ocisConfig:"backend"`
	Fallback FallbackBackend `ocisConfig:"fallback"`

	RoleBundleUUID string `ocisConfig:"role_bundle_uuid" env:"GLAUTH_ROLE_BUNDLE_ID"`

	Context context.Context
}

// Backend defined the available backend configuration.
type Backend struct {
	Datastore   string   `ocisConfig:"datastore"`
	BaseDN      string   `ocisConfig:"base_dn"`
	Insecure    bool     `ocisConfig:"insecure"`
	NameFormat  string   `ocisConfig:"name_format"`
	GroupFormat string   `ocisConfig:"group_format"`
	Servers     []string `ocisConfig:"servers"`
	SSHKeyAttr  string   `ocisConfig:"ssh_key_attr"`
	UseGraphAPI bool     `ocisConfig:"use_graph_api"`
}

// FallbackBackend defined the available fallback backend configuration.
type FallbackBackend struct {
	Datastore   string   `ocisConfig:"datastore"`
	BaseDN      string   `ocisConfig:"base_dn"`
	Insecure    bool     `ocisConfig:"insecure"`
	NameFormat  string   `ocisConfig:"name_format"`
	GroupFormat string   `ocisConfig:"group_format"`
	Servers     []string `ocisConfig:"servers"`
	SSHKeyAttr  string   `ocisConfig:"ssh_key_attr"`
	UseGraphAPI bool     `ocisConfig:"use_graph_api"`
}
