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

// Server configures a server.
type Server struct {
	Name      string
	Namespace string
	Address   string
}

// Log defines the available logging configuration.
type Log struct {
	Level  string
	Pretty bool
	Color  bool
}

// Config merges all Account config parameters.
type Config struct {
	LDAP   LDAP
	Server Server
	Log    Log
}

// New returns a new config.
func New() *Config {
	return &Config{}
}
