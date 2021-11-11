package config

import (
	"context"
	"path"

	"github.com/owncloud/ocis/ocis-pkg/shared"

	"github.com/owncloud/ocis/ocis-pkg/config/defaults"
)

// Debug defines the available debug configuration.
type Debug struct {
	Addr   string `mapstructure:"addr"`
	Token  string `mapstructure:"token"`
	Pprof  bool   `mapstructure:"pprof"`
	Zpages bool   `mapstructure:"zpages"`
}

// HTTP defines the available http configuration.
type HTTP struct {
	Addr      string `mapstructure:"addr"`
	Namespace string `mapstructure:"namespace"`
	Root      string `mapstructure:"root"`
}

// Tracing defines the available tracing configuration.
type Tracing struct {
	Enabled   bool   `mapstructure:"enabled"`
	Type      string `mapstructure:"type"`
	Endpoint  string `mapstructure:"endpoint"`
	Collector string `mapstructure:"collector"`
	Service   string `mapstructure:"service"`
}

// Ldap defined the available LDAP configuration.
type Ldap struct {
	Enabled bool   `mapstructure:"enabled"`
	Addr    string `mapstructure:"addr"`
}

// Ldaps defined the available LDAPS configuration.
type Ldaps struct {
	Addr    string `mapstructure:"addr"`
	Enabled bool   `mapstructure:"enabled"`
	Cert    string `mapstructure:"cert"`
	Key     string `mapstructure:"key"`
}

// Backend defined the available backend configuration.
type Backend struct {
	Datastore   string   `mapstructure:"datastore"`
	BaseDN      string   `mapstructure:"base_dn"`
	Insecure    bool     `mapstructure:"insecure"`
	NameFormat  string   `mapstructure:"name_format"`
	GroupFormat string   `mapstructure:"group_format"`
	Servers     []string `mapstructure:"servers"`
	SSHKeyAttr  string   `mapstructure:"ssh_key_attr"`
	UseGraphAPI bool     `mapstructure:"use_graph_api"`
}

// Config combines all available configuration parts.
type Config struct {
	*shared.Commons

	File           string      `mapstructure:"file"`
	Log            *shared.Log `mapstructure:"log"`
	Debug          Debug       `mapstructure:"debug"`
	HTTP           HTTP        `mapstructure:"http"`
	Tracing        Tracing     `mapstructure:"tracing"`
	Ldap           Ldap        `mapstructure:"ldap"`
	Ldaps          Ldaps       `mapstructure:"ldaps"`
	Backend        Backend     `mapstructure:"backend"`
	Fallback       Backend     `mapstructure:"fallback"`
	Version        string      `mapstructure:"version"`
	RoleBundleUUID string      `mapstructure:"role_bundle_uuid"`

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
		HTTP: HTTP{},
		Tracing: Tracing{
			Type:    "jaeger",
			Service: "glauth",
		},
		Ldap: Ldap{
			Enabled: true,
			Addr:    "127.0.0.1:9125",
		},
		Ldaps: Ldaps{
			Addr:    "127.0.0.1:9126",
			Enabled: true,
			Cert:    path.Join(defaults.BaseDataPath(), "ldap", "ldap.crt"),
			Key:     path.Join(defaults.BaseDataPath(), "ldap", "ldap.key"),
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
		Fallback: Backend{
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
