package config

import (
	"context"
)

// Debug defines the available debug configuration.
type Debug struct {
	Addr   string `ocisConfig:"addr" env:"GRAPH_EXPLORER_DEBUG_ADDR"`
	Token  string `ocisConfig:"token" env:"GRAPH_EXPLORER_DEBUG_TOKEN"`
	Pprof  bool   `ocisConfig:"pprof" env:"GRAPH_EXPLORER_DEBUG_PPROF"`
	Zpages bool   `ocisConfig:"zpages" env:"GRAPH_EXPLORER_DEBUG_ZPAGES"`
}

// HTTP defines the available http configuration.
type HTTP struct {
	Addr      string `ocisConfig:"addr" env:"GRAPH_EXPLORER_HTTP_ADDR"`
	Root      string `ocisConfig:"root" env:"GRAPH_EXPLORER_HTTP_ROOT"`
	Namespace string
}

// Service defines the available service configuration.
type Service struct {
	Name    string
	Version string
}

// Tracing defines the available tracing configuration.
type Tracing struct {
	Enabled   bool   `ocisConfig:"enabled" env:"OCIS_TRACING_ENABLED;GRAPH_EXPLORER_TRACING_ENABLED"`
	Type      string `ocisConfig:"type" env:"OCIS_TRACING_TYPE;GRAPH_EXPLORER_TRACING_TYPE"`
	Endpoint  string `ocisConfig:"endpoint" env:"OCIS_TRACING_ENDPOINT;GRAPH_EXPLORER_TRACING_ENDPOINT"`
	Collector string `ocisConfig:"collector" env:"OCIS_TRACING_COLLECTOR;GRAPH_EXPLORER_TRACING_COLLECTOR"`
	Service   string `ocisConfig:"service" env:"GRAPH_EXPLORER_TRACING_SERVICE"` //TODO: should this be an ID? or the same as Service.Name?
}

// Log defines the available log configuration.
type Log struct {
	Level  string `mapstructure:"level" env:"OCIS_LOG_LEVEL;GRAPH_EXPLORER_LOG_LEVEL"`
	Pretty bool   `mapstructure:"pretty" env:"OCIS_LOG_PRETTY;GRAPH_EXPLORER_LOG_PRETTY"`
	Color  bool   `mapstructure:"color" env:"OCIS_LOG_COLOR;GRAPH_EXPLORER_LOG_COLOR"`
	File   string `mapstructure:"file" env:"OCIS_LOG_FILE;GRAPH_EXPLORER_LOG_FILE"`
}

// GraphExplorer defines the available graph-explorer configuration.
type GraphExplorer struct {
	ClientID     string `ocisConfig:"client_id" env:"GRAPH_EXPLORER_CLIENT_ID"`
	Issuer       string `ocisConfig:"issuer" env:"OCIS_URL;GRAPH_EXPLORER_ISSUER"`
	GraphURLBase string `ocisConfig:"graph_url_base" env:"OCIS_URL;GRAPH_EXPLORER_GRAPH_URL_BASE"`
	GraphURLPath string `ocisConfig:"graph_url_path" env:"GRAPH_EXPLORER_GRAPH_URL_PATH"`
}

// Config combines all available configuration parts.
type Config struct {
	Service Service `ocisConfig:"service"`

	Tracing Tracing `ocisConfig:"tracing"`
	Log     Log     `ocisConfig:"log"`
	Debug   Debug   `ocisConfig:"debug"`

	HTTP HTTP `ocisConfig:"http"`

	GraphExplorer GraphExplorer `ocisConfig:"graph_explorer"`

	Context    context.Context
	Supervised bool
}

// DefaultConfig provides with a working version of a config.
func DefaultConfig() *Config {
	return &Config{
		Debug: Debug{
			Addr:   "127.0.0.1:9136",
			Token:  "",
			Pprof:  false,
			Zpages: false,
		},
		HTTP: HTTP{
			Addr:      "127.0.0.1:9135",
			Root:      "/graph-explorer",
			Namespace: "com.owncloud.web",
		},
		Service: Service{
			Name: "graph-explorer",
		},
		Tracing: Tracing{
			Type:      "jaeger",
			Endpoint:  "",
			Collector: "",
			Service:   "graph-explorer",
		},
		GraphExplorer: GraphExplorer{
			ClientID:     "ocis-explorer.js",
			Issuer:       "https://localhost:9200",
			GraphURLBase: "https://localhost:9200",
			GraphURLPath: "/graph",
		},
	}
}
