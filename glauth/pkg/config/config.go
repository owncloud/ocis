package config

import (
	"context"
)

// Config combines all available configuration parts.
type Config struct {
	Service Service

	Tracing Tracing `ocisConfig:"tracing"`
	Log     Log     `ocisConfig:"log"`
	Debug   Debug   `ocisConfig:"debug"`

	Ldap  Ldap  `ocisConfig:"ldap"`
	Ldaps Ldaps `ocisConfig:"ldaps"`

	Backend  Backend         `ocisConfig:"backend"`
	Fallback FallbackBackend `ocisConfig:"fallback"`

	RoleBundleUUID string `ocisConfig:"role_bundle_uuid" env:"GLAUTH_ROLE_BUNDLE_ID"`

	Context    context.Context
	Supervised bool
}

// Backend defined the available backend configuration.
type Backend struct {
	Datastore   string   `ocisConfig:"datastore" env:"GLAUTH_BACKEND_DATASTORE"`
	BaseDN      string   `ocisConfig:"base_dn" env:"GLAUTH_BACKEND_BASEDN"`
	Insecure    bool     `ocisConfig:"insecure" env:"GLAUTH_BACKEND_INSECURE"`
	NameFormat  string   `ocisConfig:"name_format" env:"GLAUTH_BACKEND_NAME_FORMAT"`
	GroupFormat string   `ocisConfig:"group_format" env:"GLAUTH_BACKEND_GROUP_FORMAT"`
	Servers     []string `ocisConfig:"servers"` //TODO: how to configure this via env?
	SSHKeyAttr  string   `ocisConfig:"ssh_key_attr" env:"GLAUTH_BACKEND_SSH_KEY_ATTR"`
	UseGraphAPI bool     `ocisConfig:"use_graph_api" env:"GLAUTH_BACKEND_USE_GRAPHAPI"`
}

// FallbackBackend defined the available fallback backend configuration.
type FallbackBackend struct {
	Datastore   string   `ocisConfig:"datastore" env:"GLAUTH_FALLBACK_DATASTORE"`
	BaseDN      string   `ocisConfig:"base_dn" env:"GLAUTH_FALLBACK_BASEDN"`
	Insecure    bool     `ocisConfig:"insecure" env:"GLAUTH_FALLBACK_INSECURE"`
	NameFormat  string   `ocisConfig:"name_format" env:"GLAUTH_FALLBACK_NAME_FORMAT"`
	GroupFormat string   `ocisConfig:"group_format" env:"GLAUTH_FALLBACK_GROUP_FORMAT"`
	Servers     []string `ocisConfig:"servers"` //TODO: how to configure this via env?
	SSHKeyAttr  string   `ocisConfig:"ssh_key_attr" env:"GLAUTH_FALLBACK_SSH_KEY_ATTR"`
	UseGraphAPI bool     `ocisConfig:"use_graph_api" env:"GLAUTH_FALLBACK_USE_GRAPHAPI"`
}
