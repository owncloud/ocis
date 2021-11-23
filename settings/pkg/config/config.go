package config

import (
	"context"
	"path"

	"github.com/owncloud/ocis/ocis-pkg/shared"

	"github.com/owncloud/ocis/ocis-pkg/config/defaults"
)

// Debug defines the available debug configuration.
type Debug struct {
	Addr   string `ocisConfig:"addr"`
	Token  string `ocisConfig:"token"`
	Pprof  bool   `ocisConfig:"pprof"`
	Zpages bool   `ocisConfig:"zpages"`
}

// CORS defines the available cors configuration.
type CORS struct {
	AllowedOrigins   []string `ocisConfig:"allowed_origins"`
	AllowedMethods   []string `ocisConfig:"allowed_methods"`
	AllowedHeaders   []string `ocisConfig:"allowed_headers"`
	AllowCredentials bool     `ocisConfig:"allow_credentials"`
}

// HTTP defines the available http configuration.
type HTTP struct {
	Addr      string `ocisConfig:"addr"`
	Namespace string `ocisConfig:"namespace"`
	Root      string `ocisConfig:"root"`
	CacheTTL  int    `ocisConfig:"cache_ttl"`
	CORS      CORS   `ocisConfig:"cors"`
}

// GRPC defines the available grpc configuration.
type GRPC struct {
	Addr      string `ocisConfig:"grpc"`
	Namespace string `ocisConfig:"namespace"`
}

// Service provides configuration options for the service
type Service struct {
	Name     string `ocisConfig:"name"`
	Version  string `ocisConfig:"version"`
	DataPath string `ocisConfig:"data_path"`
}

// Tracing defines the available tracing configuration.
type Tracing struct {
	Enabled   bool   `ocisConfig:"enabled"`
	Type      string `ocisConfig:"type"`
	Endpoint  string `ocisConfig:"endpoint"`
	Collector string `ocisConfig:"collector"`
	Service   string `ocisConfig:"service"`
}

// Asset undocumented
type Asset struct {
	Path string `ocisConfig:"asset"`
}

// TokenManager is the config for using the reva token manager
type TokenManager struct {
	JWTSecret string `ocisConfig:"jwt_secret"`
}

// Config combines all available configuration parts.
type Config struct {
	*shared.Commons

	File         string       `ocisConfig:"file"`
	Service      Service      `ocisConfig:"service"`
	Log          *shared.Log  `ocisConfig:"log"`
	Debug        Debug        `ocisConfig:"debug"`
	HTTP         HTTP         `ocisConfig:"http"`
	GRPC         GRPC         `ocisConfig:"grpc"`
	Tracing      Tracing      `ocisConfig:"tracing"`
	Asset        Asset        `ocisConfig:"asset"`
	TokenManager TokenManager `ocisConfig:"token_manager"`

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
