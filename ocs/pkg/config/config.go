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

// CORS defines the available cors configuration.
type CORS struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	AllowCredentials bool
}

// HTTP defines the available http configuration.
type HTTP struct {
	Addr string
	Root string
	CORS CORS
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

// Reva defines all available REVA configuration.
type Reva struct {
	Address string
}

// TokenManager is the config for using the reva token manager
type TokenManager struct {
	JWTSecret string
}

// IdentityManagement keeps track of the OIDC address. This is because Reva requisite of uniqueness for users
// is based in the combination of IDP hostname + UserID. For more information see:
// https://github.com/cs3org/reva/blob/4fd0229f13fae5bc9684556a82dbbd0eced65ef9/pkg/storage/utils/decomposedfs/node/node.go#L856-L865
type IdentityManagement struct {
	Address string
}

// Config combines all available configuration parts.
type Config struct {
	File               string
	Log                Log
	Debug              Debug
	HTTP               HTTP
	Tracing            Tracing
	TokenManager       TokenManager
	Service            Service
	AccountBackend     string
	Reva               Reva
	StorageUsersDriver string
	MachineAuthAPIKey  string
	IdentityManagement IdentityManagement

	Context    context.Context
	Supervised bool
}

// New initializes a new configuration with or without defaults.
func New() *Config {
	return &Config{}
}
