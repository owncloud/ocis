package config

import (
	"context"
	"path"

	"github.com/owncloud/ocis/ocis-pkg/shared"

	"github.com/owncloud/ocis/ocis-pkg/config/defaults"
)

// Debug defines the available debug configuration.
type Debug struct {
	Addr   string `ocisConfig:"addr" env:"SETTINGS_DEBUG_ADDR"`
	Token  string `ocisConfig:"token" env:"SETTINGS_DEBUG_TOKEN"`
	Pprof  bool   `ocisConfig:"pprof" env:"SETTINGS_DEBUG_PPROF"`
	Zpages bool   `ocisConfig:"zpages" env:"SETTINGS_DEBUG_ZPAGES"`
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
	Addr      string `ocisConfig:"addr" env:"SETTINGS_HTTP_ADDR"`
	Namespace string
	Root      string `ocisConfig:"root" env:"SETTINGS_HTTP_ROOT"`
	CacheTTL  int    `ocisConfig:"cache_ttl" env:"SETTINGS_CACHE_TTL"`
	CORS      CORS   `ocisConfig:"cors"`
}

// GRPC defines the available grpc configuration.
type GRPC struct {
	Addr      string `ocisConfig:"addr" env:"SETTINGS_GRPC_ADDR"`
	Namespace string
}

// Service defines the available service configuration.
type Service struct {
	Name    string
	Version string
}

// Tracing defines the available tracing configuration.
type Tracing struct {
	Enabled   bool   `ocisConfig:"enabled" env:"OCIS_TRACING_ENABLED;SETTINGS_TRACING_ENABLED"`
	Type      string `ocisConfig:"type" env:"OCIS_TRACING_TYPE;SETTINGS_TRACING_TYPE"`
	Endpoint  string `ocisConfig:"endpoint" env:"OCIS_TRACING_ENDPOINT;SETTINGS_TRACING_ENDPOINT"`
	Collector string `ocisConfig:"collector" env:"OCIS_TRACING_COLLECTOR;SETTINGS_TRACING_COLLECTOR"`
	Service   string `ocisConfig:"service" env:"SETTINGS_TRACING_SERVICE"` //TODO: should this be an ID? or the same as Service.Name?
}

// Log defines the available log configuration.
type Log struct {
	Level  string `mapstructure:"level" env:"OCIS_LOG_LEVEL;SETTINGS_LOG_LEVEL"`
	Pretty bool   `mapstructure:"pretty" env:"OCIS_LOG_PRETTY;SETTINGS_LOG_PRETTY"`
	Color  bool   `mapstructure:"color" env:"OCIS_LOG_COLOR;SETTINGS_LOG_COLOR"`
	File   string `mapstructure:"file" env:"OCIS_LOG_FILE;SETTINGS_LOG_FILE"`
}

// Asset defines the available asset configuration.
type Asset struct {
	Path string `ocisConfig:"path" env:"SETTINGS_ASSET_PATH"`
}

// TokenManager is the config for using the reva token manager
type TokenManager struct {
	JWTSecret string `ocisConfig:"jwt_secret" env:"OCIS_JWT_SECRET;SETTINGS_JWT_SECRET"`
}

// Config combines all available configuration parts.
type Config struct {
	*shared.Commons

	Service Service `ocisConfig:"service"`

	Tracing Tracing `ocisConfig:"tracing"`
	Log     Log     `ocisConfig:"log"`
	Debug   Debug   `ocisConfig:"debug"`

	HTTP HTTP `ocisConfig:"http"`
	GRPC GRPC `ocisConfig:"grpc"`

	DataPath     string       `ocisConfig:"data_path" env:"SETTINGS_DATA_PATH"`
	Asset        Asset        `ocisConfig:"asset"`
	TokenManager TokenManager `ocisConfig:"token_manager"`

	Context    context.Context
	Supervised bool
}

// DefaultConfig provides sane bootstrapping defaults.
func DefaultConfig() *Config {
	return &Config{
		Service: Service{
			Name: "settings",
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
		DataPath: path.Join(defaults.BaseDataPath(), "settings"),
		Asset: Asset{
			Path: "",
		},
		TokenManager: TokenManager{
			JWTSecret: "Pive-Fumkiu4",
		},
	}
}
