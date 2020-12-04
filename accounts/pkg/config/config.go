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
	CacheTTL  int
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
	HashDifficulty   int
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

// Repo defines which storage implementation is to be used.
type Repo struct {
	Disk Disk
	CS3  CS3
}

// Disk is the local disk implementation of the storage.
type Disk struct {
	Path string
}

// CS3 is the cs3 implementation of the storage.
type CS3 struct {
	ProviderAddr string
	DataURL      string
	DataPrefix   string
	JWTSecret    string
}

// ServiceUser defines the user required for EOS.
type ServiceUser struct {
	UUID     string
	Username string
	UID      int64
	GID      int64
}

// Index defines config for indexes.
type Index struct {
	UID, GID Bound
}

// Bound defines a lower and upper bound.
type Bound struct {
	Lower, Upper int64
}

// Tracing defines the available tracing configuration.
type Tracing struct {
	Enabled   bool
	Type      string
	Endpoint  string
	Collector string
	Service   string
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
	Repo         Repo
	Index        Index
	ServiceUser  ServiceUser
	Tracing      Tracing
}

// New returns a new config.
func New() *Config {
	return &Config{}
}
