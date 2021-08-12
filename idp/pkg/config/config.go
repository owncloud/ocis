package config

import (
	"context"

	"github.com/libregraph/lico/bootstrap"
)

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
	Addr    string
	Root    string
	TLSCert string
	TLSKey  string
	TLS     bool
}

// Ldap defines the available LDAP configuration.
type Ldap struct {
	URI               string
	BindDN            string
	BindPassword      string
	BaseDN            string
	Scope             string
	LoginAttribute    string
	EmailAttribute    string
	NameAttribute     string
	UUIDAttribute     string
	UUIDAttributeType string
	Filter            string
}

// Service defines the available service configuration.
type Service struct {
	Name      string
	Namespace string
	Version   string
}

// Tracing defines the available tracing configuration.
type Tracing struct {
	Enabled   bool
	Type      string
	Endpoint  string
	Collector string
	Service   string
}

// Asset defines the available asset configuration.
type Asset struct {
	Path string
}

// Config combines all available configuration parts.
type Config struct {
	File    string
	Log     Log
	Debug   Debug
	HTTP    HTTP
	Tracing Tracing
	Asset   Asset
	IDP     bootstrap.Config
	Ldap    Ldap
	Service Service

	Context    context.Context
	Supervised bool
}

// New initializes a new configuration with or without defaults.
func New() *Config {
	return &Config{}
}
