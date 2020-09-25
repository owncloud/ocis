// Package config should be moved to internal
package config

// LDAP defines the available ldap configuration.
type LDAP struct {
	Hostname     string
	Port         int
	BaseDN       string
	UserFilter   string
	GroupFilter  string
	BindDN       string
	BindPassword string
	IDP          string
	Schema       LDAPSchema
}

// LDAPSchema defines the available ldap schema configuration.
type LDAPSchema struct {
	AccountID   string
	Identities  string
	Username    string
	DisplayName string
	Mail        string
	Groups      string
}

// HTTP defines the available http configuration.
type HTTP struct {
	Addr      string
	Namespace string
	Root      string
}

// GRPC defines the available grpc configuration.
type GRPC struct {
	Addr      string
	Namespace string
}

// Server configures a server.
type Server struct {
	Version          string
	Name             string
	AccountsDataPath string
}

// Asset defines the available asset configuration.
type Asset struct {
	Path string
}

// TokenManager is the config for using the reva token manager
type TokenManager struct {
	JWTSecret string
}

// Log defines the available logging configuration.
type Log struct {
	Level  string
	Pretty bool
	Color  bool
}

// Config merges all Account config parameters.
type Config struct {
	LDAP         LDAP
	HTTP         HTTP
	GRPC         GRPC
	Server       Server
	Asset        Asset
	Log          Log
	TokenManager TokenManager
}

// New returns a new config.
func New() *Config {
	return &Config{}
}
