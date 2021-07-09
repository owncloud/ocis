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

// Server configures a server.
type Server struct {
	Version string
	Name    string
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
	Network      string
	Address      string
	UserName     string
	Password     string
	BaseDNUsers  string
	BaseDNGroups string
}

// OpenIDConnect defined the available OpenID Connect configuration.
type OpenIDConnect struct {
	Endpoint    string
	Realm       string
	SigningAlgs []string
	Insecure    bool
}

// Reva defines all available REVA configuration.
type Reva struct {
	Address string
}

// TokenManager is the config for using the reva token manager
type TokenManager struct {
	JWTSecret string
}

type Spaces struct {
	WebDavBase string
}

// Config combines all available configuration parts.
type Config struct {
	File            string
	WebdavNamespace string
	Log             Log
	Debug           Debug
	HTTP            HTTP
	Server          Server
	Tracing         Tracing
	Ldap            Ldap
	OpenIDConnect   OpenIDConnect
	Reva            Reva
	TokenManager    TokenManager
	Spaces          Spaces

	Context    context.Context
	Supervised bool
}

// New initializes a new configuration with or without defaults.
func New() *Config {
	return &Config{}
}
