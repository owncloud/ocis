package config

import (
	"context"
	"path"

	"github.com/owncloud/ocis/ocis-pkg/shared"

	"github.com/owncloud/ocis/ocis-pkg/config/defaults"
)

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
	Addr      string
	Namespace string
	Root      string
	CacheTTL  int
	CORS      CORS
}

// GRPC defines the available grpc configuration.
type GRPC struct {
	Addr      string
	Namespace string
}

// Service provides configuration options for the service
type Service struct {
	Name     string
	Version  string
	DataPath string
}

// Tracing defines the available tracing configuration.
type Tracing struct {
	Enabled   bool
	Type      string
	Endpoint  string
	Collector string
	Service   string
}

// Asset undocumented
type Asset struct {
	Path string
}

// TokenManager is the config for using the reva token manager
type TokenManager struct {
	JWTSecret string
}

// Config combines all available configuration parts.
type Config struct {
	File         string
	Service      Service
	Log          shared.Log
	Debug        Debug
	HTTP         HTTP
	GRPC         GRPC
	Tracing      Tracing
	Asset        Asset
	TokenManager TokenManager

	Context    context.Context
	Supervised bool
}

// New initializes a new configuration with or without defaults.
func New() *Config {
	return &Config{}
}

// DefaultConfig provides sane bootstrapping defaults.
func DefaultConfig() *Config {
	return &Config{
		Log: shared.Log{},
		Service: Service{
			Name:     "settings",
			DataPath: path.Join(defaults.BaseDataPath(), "settings"),
		},
		Debug: Debug{
			Addr:   "127.0.0.1:9194",
			Token:  "",
			Pprof:  false,
			Zpages: false,
		},
		HTTP: HTTP{
			Addr:      "127.0.0.1:9190",
			Namespace: "com.owncloud.web",
			Root:      "/",
			CacheTTL:  604800, // 7 days
			CORS: CORS{
				AllowedOrigins:   []string{"*"},
				AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
				AllowedHeaders:   []string{"Authorization", "Origin", "Content-Type", "Accept", "X-Requested-With"},
				AllowCredentials: true,
			},
		},
		GRPC: GRPC{
			Addr:      "127.0.0.1:9191",
			Namespace: "com.owncloud.api",
		},
		Tracing: Tracing{
			Enabled:   false,
			Type:      "jaeger",
			Endpoint:  "",
			Collector: "",
			Service:   "settings",
		},
		Asset: Asset{
			Path: "",
		},
		TokenManager: TokenManager{
			JWTSecret: "Pive-Fumkiu4",
		},
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
