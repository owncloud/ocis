package config

import (
	"context"

	"github.com/owncloud/ocis/ocis-pkg/shared"
)

// Debug defines the available debug configuration.
type Debug struct {
	Addr   string `ocisConfig:"addr" env:"WEBDAV_DEBUG_ADDR"`
	Token  string `ocisConfig:"token" env:"WEBDAV_DEBUG_TOKEN"`
	Pprof  bool   `ocisConfig:"pprof" env:"WEBDAV_DEBUG_PPROF"`
	Zpages bool   `ocisConfig:"zpages" env:"WEBDAV_DEBUG_ZPAGES"`
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
	Addr      string `ocisConfig:"addr" env:"WEBDAV_HTTP_ADDR"`
	Namespace string
	Root      string `ocisConfig:"root" env:"WEBDAV_HTTP_ROOT"`
	CORS      CORS   `ocisConfig:"cors"`
}

// Service defines the available service configuration.
type Service struct {
	Name    string
	Version string
}

// Tracing defines the available tracing configuration.
type Tracing struct {
	Enabled   bool   `ocisConfig:"enabled" env:"OCIS_TRACING_ENABLED;WEBDAV_TRACING_ENABLED"`
	Type      string `ocisConfig:"type" env:"OCIS_TRACING_TYPE;WEBDAV_TRACING_TYPE"`
	Endpoint  string `ocisConfig:"endpoint" env:"OCIS_TRACING_ENDPOINT;WEBDAV_TRACING_ENDPOINT"`
	Collector string `ocisConfig:"collector" env:"OCIS_TRACING_COLLECTOR;WEBDAV_TRACING_COLLECTOR"`
	Service   string `ocisConfig:"service" env:"WEBDAV_TRACING_SERVICE"` //TODO: should this be an ID? or the same as Service.Name?
}

// Log defines the available log configuration.
type Log struct {
	Level  string `mapstructure:"level" env:"OCIS_LOG_LEVEL;WEBDAV_LOG_LEVEL"`
	Pretty bool   `mapstructure:"pretty" env:"OCIS_LOG_PRETTY;WEBDAV_LOG_PRETTY"`
	Color  bool   `mapstructure:"color" env:"OCIS_LOG_COLOR;WEBDAV_LOG_COLOR"`
	File   string `mapstructure:"file" env:"OCIS_LOG_FILE;WEBDAV_LOG_FILE"`
}

// Config combines all available configuration parts.
type Config struct {
	*shared.Commons

	Service Service `ocisConfig:"service"`

	Tracing Tracing `ocisConfig:"tracing"`
	Log     Log     `ocisConfig:"log"`
	Debug   Debug   `ocisConfig:"debug"`

	HTTP HTTP `ocisConfig:"http"`

	OcisPublicURL   string `ocisConfig:"ocis_public_url" env:"OCIS_URL;OCIS_PUBLIC_URL"`
	WebdavNamespace string `ocisConfig:"webdav_namespace" env:"STORAGE_WEBDAV_NAMESPACE"` //TODO: prevent this cross config

	Context    context.Context
	Supervised bool
}

func DefaultConfig() *Config {
	return &Config{
		Debug: Debug{
			Addr:   "127.0.0.1:9119",
			Token:  "",
			Pprof:  false,
			Zpages: false,
		},
		HTTP: HTTP{
			Addr:      "127.0.0.1:9115",
			Root:      "/",
			Namespace: "com.owncloud.web",
			CORS: CORS{
				AllowedOrigins:   []string{"*"},
				AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
				AllowedHeaders:   []string{"Authorization", "Origin", "Content-Type", "Accept", "X-Requested-With"},
				AllowCredentials: true,
			},
		},
		Service: Service{
			Name: "webdav",
		},
		Tracing: Tracing{
			Enabled:   false,
			Type:      "jaeger",
			Endpoint:  "",
			Collector: "",
			Service:   "webdav",
		},
		OcisPublicURL:   "https://127.0.0.1:9200",
		WebdavNamespace: "/home",
	}
}
