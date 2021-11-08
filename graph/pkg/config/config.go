package config

import (
	"context"

	"github.com/owncloud/ocis/ocis-pkg/shared"
)

// Debug defines the available debug configuration.
type Debug struct {
	Addr   string `mapstructure:"addr"`
	Token  string `mapstructure:"token"`
	Pprof  bool   `mapstructure:"pprof"`
	Zpages bool   `mapstructure:"zpages"`
}

// HTTP defines the available http configuration.
type HTTP struct {
	Addr      string `mapstructure:"addr"`
	Namespace string `mapstructure:"namespace"`
	Root      string `mapstructure:"root"`
}

// Server configures a server.
type Server struct {
	Version string `mapstructure:"version"`
	Name    string `mapstructure:"name"`
}

// Tracing defines the available tracing configuration.
type Tracing struct {
	Enabled   bool   `mapstructure:"enabled"`
	Type      string `mapstructure:"type"`
	Endpoint  string `mapstructure:"endpoint"`
	Collector string `mapstructure:"collector"`
	Service   string `mapstructure:"service"`
}

// Reva defines all available REVA configuration.
type Reva struct {
	Address string `mapstructure:"address"`
}

// TokenManager is the config for using the reva token manager
type TokenManager struct {
	JWTSecret string `mapstructure:"jwt_secret"`
}

type Spaces struct {
	WebDavBase   string `mapstructure:"webdav_base"`
	WebDavPath   string `mapstructure:"webdav_path"`
	DefaultQuota string `mapstructure:"default_quota"`
}

// Config combines all available configuration parts.
type Config struct {
	File         string       `mapstructure:"file"`
	Log          shared.Log   `mapstructure:"log"`
	Debug        Debug        `mapstructure:"debug"`
	HTTP         HTTP         `mapstructure:"http"`
	Server       Server       `mapstructure:"server"`
	Tracing      Tracing      `mapstructure:"tracing"`
	Reva         Reva         `mapstructure:"reva"`
	TokenManager TokenManager `mapstructure:"token_manager"`
	Spaces       Spaces       `mapstructure:"spaces"`

	Context    context.Context
	Supervised bool
}

// New initializes a new configuration with or without defaults.
func New() *Config {
	return &Config{}
}

func DefaultConfig() *Config {
	return &Config{
		Log: shared.Log{},
		Debug: Debug{
			Addr:  "127.0.0.1:9124",
			Token: "",
		},
		HTTP: HTTP{
			Addr:      "127.0.0.1:9120",
			Namespace: "com.owncloud.web",
			Root:      "/graph",
		},
		Server: Server{},
		Tracing: Tracing{
			Enabled: false,
			Type:    "jaeger",
			Service: "graph",
		},
		Reva: Reva{
			Address: "127.0.0.1:9142",
		},
		TokenManager: TokenManager{
			JWTSecret: "Pive-Fumkiu4",
		},
		Spaces: Spaces{
			WebDavBase:   "https://localhost:9200",
			WebDavPath:   "/dav/spaces/",
			DefaultQuota: "1000000000",
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
