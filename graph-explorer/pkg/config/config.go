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

// HTTP defines the available http configuration.
type HTTP struct {
	Addr      string `ocisConfig:"addr"`
	Root      string `ocisConfig:"root"`
	Namespace string `ocisConfig:"namespace"`
}

// Server configures a server.
type Server struct {
	Version string `ocisConfig:"version"`
	Name    string `ocisConfig:"name"`
}

// Tracing defines the available tracing configuration.
type Tracing struct {
	Enabled   bool   `ocisConfig:"enabled"`
	Type      string `ocisConfig:"type"`
	Endpoint  string `ocisConfig:"endpoint"`
	Collector string `ocisConfig:"collector"`
	Service   string `ocisConfig:"service"`
}

// GraphExplorer defines the available graph-explorer configuration.
type GraphExplorer struct {
	ClientID     string `ocisConfig:"client_id"`
	Issuer       string `ocisConfig:"issuer"`
	GraphURLBase string `ocisConfig:"graph_url_base"`
	GraphURLPath string `ocisConfig:"graph_url_path"`
}

// Config combines all available configuration parts.
type Config struct {
	File          string        `ocisConfig:"file"`
	Log           shared.Log    `ocisConfig:"log"`
	Debug         Debug         `ocisConfig:"debug"`
	HTTP          HTTP          `ocisConfig:"http"`
	Server        Server        `ocisConfig:"server"`
	Tracing       Tracing       `ocisConfig:"tracing"`
	GraphExplorer GraphExplorer `ocisConfig:"graph_explorer"`

	Context    context.Context
	Supervised bool
}

// New initializes a new configuration with or without defaults.
func New() *Config {
	return &Config{}
}

// DefaultConfig provides with a working version of a config.
func DefaultConfig() *Config {
	return &Config{
		Log: shared.Log{},
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
		Server: Server{},
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

// GetEnv fetches a list of known env variables for this extension. It is to be used by gookit, as it provides a list
// with all the environment variables an extension supports.
func GetEnv() []string {
	var r = make([]string, len(structMappings(&Config{})))
	for i := range structMappings(&Config{}) {
		r = append(r, structMappings(&Config{})[i].EnvVars...)
	}

	return r
}
