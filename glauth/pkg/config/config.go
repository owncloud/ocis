package config

import (
	"context"
	"path"

	"github.com/owncloud/ocis/ocis-pkg/shared"

	"github.com/owncloud/ocis/ocis-pkg/config/defaults"
)

// Debug defines the available debug configuration.
type Debug struct {
	Addr   string `ocisConfig:"addr" env:"GLAUTH_DEBUG_ADDR"`
	Token  string `ocisConfig:"token" env:"GLAUTH_DEBUG_TOKEN"`
	Pprof  bool   `ocisConfig:"pprof" env:"GLAUTH_DEBUG_PPROF"`
	Zpages bool   `ocisConfig:"zpages" env:"GLAUTH_DEBUG_ZPAGES"`
}

// Service defines the available service configuration.
type Service struct {
	Name    string
	Version string
}

// Tracing defines the available tracing configuration.
type Tracing struct {
	Enabled   bool   `ocisConfig:"enabled" env:"OCIS_TRACING_ENABLED;GLAUTH_TRACING_ENABLED"`
	Type      string `ocisConfig:"type" env:"OCIS_TRACING_TYPE;GLAUTH_TRACING_TYPE"`
	Endpoint  string `ocisConfig:"endpoint" env:"OCIS_TRACING_ENDPOINT;GLAUTH_TRACING_ENDPOINT"`
	Collector string `ocisConfig:"collector" env:"OCIS_TRACING_COLLECTOR;GLAUTH_TRACING_COLLECTOR"`
	Service   string `ocisConfig:"service" env:"GLAUTH_TRACING_SERVICE"` // TODO:
}

// Log defines the available log configuration.
type Log struct {
	Level  string `mapstructure:"level" env:"OCIS_LOG_LEVEL;GLAUTH_LOG_LEVEL"`
	Pretty bool   `mapstructure:"pretty" env:"OCIS_LOG_PRETTY;GLAUTH_LOG_PRETTY"`
	Color  bool   `mapstructure:"color" env:"OCIS_LOG_COLOR;GLAUTH_LOG_COLOR"`
	File   string `mapstructure:"file" env:"OCIS_LOG_FILE;GLAUTH_LOG_FILE"`
}

// Ldap defined the available LDAP configuration.
type Ldap struct {
	Enabled   bool   `ocisConfig:"enabled" env:"GLAUTH_LDAP_ENABLED"`
	Addr      string `ocisConfig:"addr" env:"GLAUTH_LDAP_ADDR"`
	Namespace string
}

// Ldaps defined the available LDAPS configuration.
type Ldaps struct {
	Enabled   bool   `ocisConfig:"enabled" env:"GLAUTH_LDAPS_ENABLED"`
	Addr      string `ocisConfig:"addr" env:"GLAUTH_LDAPS_ADDR"`
	Namespace string
	Cert      string `ocisConfig:"cert" env:"GLAUTH_LDAPS_CERT"`
	Key       string `ocisConfig:"key" env:"GLAUTH_LDAPS_KEY"`
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

// Config combines all available configuration parts.
type Config struct {
	*shared.Commons

	Log            Log             `ocisConfig:"log"`
	Debug          Debug           `ocisConfig:"debug"`
	Service        Service         `ocisConfig:"service"`
	Tracing        Tracing         `ocisConfig:"tracing"`
	Ldap           Ldap            `ocisConfig:"ldap"`
	Ldaps          Ldaps           `ocisConfig:"ldaps"`
	Backend        Backend         `ocisConfig:"backend"`
	Fallback       FallbackBackend `ocisConfig:"fallback"`
	RoleBundleUUID string          `ocisConfig:"role_bundle_uuid" env:"GLAUTH_ROLE_BUNDLE_ID"`

	Context    context.Context
	Supervised bool
}

// New initializes a new configuration with or without defaults.
func New() *Config {
	return &Config{}
}

func DefaultConfig() *Config {
	return &Config{
		Debug: Debug{
			Addr: "127.0.0.1:9129",
		},
		Tracing: Tracing{
			Type:    "jaeger",
			Service: "glauth",
		},
		Service: Service{
			Name: "glauth",
		},
		Ldap: Ldap{
			Enabled:   true,
			Addr:      "127.0.0.1:9125",
			Namespace: "com.owncloud.ldap",
		},
		Ldaps: Ldaps{
			Enabled:   true,
			Addr:      "127.0.0.1:9126",
			Namespace: "com.owncloud.ldaps",
			Cert:      path.Join(defaults.BaseDataPath(), "ldap", "ldap.crt"),
			Key:       path.Join(defaults.BaseDataPath(), "ldap", "ldap.key"),
		},
		Backend: Backend{
			Datastore:   "accounts",
			BaseDN:      "dc=ocis,dc=test",
			Insecure:    false,
			NameFormat:  "cn",
			GroupFormat: "ou",
			Servers:     nil,
			SSHKeyAttr:  "sshPublicKey",
			UseGraphAPI: true,
		},
		Fallback: FallbackBackend{
			Datastore:   "",
			BaseDN:      "dc=ocis,dc=test",
			Insecure:    false,
			NameFormat:  "cn",
			GroupFormat: "ou",
			Servers:     nil,
			SSHKeyAttr:  "sshPublicKey",
			UseGraphAPI: true,
		},
		RoleBundleUUID: "71881883-1768-46bd-a24d-a356a2afdf7f", // BundleUUIDRoleAdmin
	}
}
