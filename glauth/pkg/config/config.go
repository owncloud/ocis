package config

import (
	"context"
	"fmt"
	"path"
	"reflect"

	gofig "github.com/gookit/config/v2"
	"github.com/owncloud/ocis/ocis-pkg/config/defaults"
)

// Log defines the available logging configuration.
type Log struct {
	Level  string `mapstructure:"level"`
	Pretty bool   `mapstructure:"pretty"`
	Color  bool   `mapstructure:"color"`
	File   string `mapstructure:"file"`
}

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
	File           string  `mapstructure:"file"`
	Log            Log     `mapstructure:"log"`
	Debug          Debug   `mapstructure:"debug"`
	HTTP           HTTP    `mapstructure:"http"`
	Tracing        Tracing `mapstructure:"tracing"`
	Ldap           Ldap    `mapstructure:"ldap"`
	Ldaps          Ldaps   `mapstructure:"ldaps"`
	Backend        Backend `mapstructure:"backend"`
	Fallback       Backend `mapstructure:"fallback"`
	Version        string  `mapstructure:"version"`
	RoleBundleUUID string  `mapstructure:"role_bundle_uuid"`

	Context    context.Context
	Supervised bool
}

// New initializes a new configuration with or without defaults.
func New() *Config {
	return &Config{}
}

func DefaultConfig() *Config {
	return &Config{
		Log: Log{},
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

// GetEnv fetches a list of known env variables for this extension. It is to be used by gookit, as it provides a list
// with all the environment variables an extension supports.
func GetEnv() []string {
	var r = make([]string, len(structMappings(&Config{})))
	for i := range structMappings(&Config{}) {
		r = append(r, structMappings(&Config{})[i].EnvVars...)
	}

	return r
}

// UnmapEnv loads values from the gooconf.Config argument and sets them in the expected destination.
func (c *Config) UnmapEnv(gooconf *gofig.Config) error {
	vals := structMappings(c)
	for i := range vals {
		for j := range vals[i].EnvVars {
			// we need to guard against v != "" because this is the condition that checks that the value is set from the environment.
			// the `ok` guard is not enough, apparently.
			if v, ok := gooconf.GetValue(vals[i].EnvVars[j]); ok && v != "" {

				// get the destination type from destination
				switch reflect.ValueOf(vals[i].Destination).Type().String() {
				case "*bool":
					r := gooconf.Bool(vals[i].EnvVars[j])
					*vals[i].Destination.(*bool) = r
				case "*string":
					r := gooconf.String(vals[i].EnvVars[j])
					*vals[i].Destination.(*string) = r
				case "*int":
					r := gooconf.Int(vals[i].EnvVars[j])
					*vals[i].Destination.(*int) = r
				case "*float64":
					// defaults to float64
					r := gooconf.Float(vals[i].EnvVars[j])
					*vals[i].Destination.(*float64) = r
				default:
					// it is unlikely we will ever get here. Let this serve more as a runtime check for when debugging.
					return fmt.Errorf("invalid type for env var: `%v`", vals[i].EnvVars[j])
				}
			}
		}
	}

	return nil
}
