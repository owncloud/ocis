package config

import "context"

// Log defines the available logging configuration.
type Log struct {
	Level  string
	Pretty bool
	Color  bool
	File   string
}

// Debug defines the available debug configuration.
type Debug struct {
	Addr   string
	Token  string
	Pprof  bool
	Zpages bool
}

// HTTP defines the available http configuration.
type HTTP struct {
	Addr      string
	Namespace string
	Root      string
}

// Tracing defines the available tracing configuration.
type Tracing struct {
	Enabled   bool
	Type      string
	Endpoint  string
	Collector string
	Service   string
}

// Ldap defined the available LDAP configuration.
type Ldap struct {
	Address string
	Enabled bool
}

// Ldaps defined the available LDAPS configuration.
type Ldaps struct {
	Ldap
	Cert string
	Key  string
}

// Backend defined the available backend configuration.
type Backend struct {
	Datastore   string
	BaseDN      string
	Insecure    bool
	NameFormat  string
	GroupFormat string
	Servers     []string
	SSHKeyAttr  string
	UseGraphAPI bool
}

// Config combines all available configuration parts.
type Config struct {
	File           string
	Log            Log
	Debug          Debug
	HTTP           HTTP
	Tracing        Tracing
	Ldap           Ldap
	Ldaps          Ldaps
	Backend        Backend
	Fallback       Backend
	Version        string
	RoleBundleUUID string

	Context    context.Context
	Supervised bool
}

// New initializes a new configuration with or without defaults.
func New() *Config {
	return &Config{}
}
