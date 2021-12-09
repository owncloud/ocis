package config

import (
	"context"

	"github.com/owncloud/ocis/ocis-pkg/shared"
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
	Addr string `ocisConfig:"addr"`
	Root string `ocisConfig:"root"`
	CORS CORS   `ocisConfig:"cors"`
}

// Service defines the available service configuration.
type Service struct {
	Name      string `ocisConfig:"name"`
	Namespace string `ocisConfig:"namespace"`
	Version   string `ocisConfig:"version"`
}

// Tracing defines the available tracing configuration.
type Tracing struct {
	Enabled   bool   `ocisConfig:"enabled"`
	Type      string `ocisConfig:"type"`
	Endpoint  string `ocisConfig:"endpoint"`
	Collector string `ocisConfig:"collector"`
	Service   string `ocisConfig:"service"`
}

// Config combines all available configuration parts.
type Config struct {
	*shared.Commons

	File            string      `ocisConfig:"file"`
	Log             *shared.Log `ocisConfig:"log"`
	Debug           Debug       `ocisConfig:"debug"`
	HTTP            HTTP        `ocisConfig:"http"`
	Tracing         Tracing     `ocisConfig:"tracing"`
	Service         Service     `ocisConfig:"service"`
	OcisPublicURL   string      `ocisConfig:"ocis_public_url"`
	WebdavNamespace string      `ocisConfig:"webdav_namespace"`

	Context    context.Context
	Supervised bool
}

// New initializes a new configuration with or without defaults.
func New() *Config {
	return &Config{}
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
			Addr: "127.0.0.1:9115",
			Root: "/",
			CORS: CORS{
				AllowedOrigins:   []string{"*"},
				AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
				AllowedHeaders:   []string{"Authorization", "Origin", "Content-Type", "Accept", "X-Requested-With"},
				AllowCredentials: true,
			},
		},
		Tracing: Tracing{
			Enabled:   false,
			Type:      "jaeger",
			Endpoint:  "",
			Collector: "",
			Service:   "webdav",
		},
		Service: Service{
			Name:      "webdav",
			Namespace: "com.owncloud.web",
		},
		OcisPublicURL:   "https://127.0.0.1:9200",
		WebdavNamespace: "/home",
	}
}
