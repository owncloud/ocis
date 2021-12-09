package config

import (
	"context"
	"path"

	"github.com/owncloud/ocis/ocis-pkg/config/defaults"
	"github.com/owncloud/ocis/ocis-pkg/shared"
)

// Debug defines the available debug configuration.
type Debug struct {
	Addr   string `ocisConfig:"addr"`
	Token  string `ocisConfig:"token"`
	Pprof  bool   `ocisConfig:"pprof"`
	Zpages bool   `ocisConfig:"zpages"`
}

// GRPC defines the available grpc configuration.
type GRPC struct {
	Addr string `ocisConfig:"addr"`
	Root string `ocisConfig:"root"`
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
	File     string     `ocisConfig:"file"`
	Log      shared.Log `ocisConfig:"log"`
	Debug    Debug      `ocisConfig:"debug"`
	GRPC     GRPC       `ocisConfig:"grpc"`
	Tracing  Tracing    `ocisConfig:"tracing"`
	Datapath string     `ocisConfig:"data_path"`
	Service  Service    `ocisConfig:"service"`

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
			Addr:   "127.0.0.1:9464",
			Token:  "",
			Pprof:  false,
			Zpages: false,
		},
		GRPC: GRPC{
			Addr: "127.0.0.1:9460",
		},
		Tracing: Tracing{
			Enabled:   false,
			Type:      "jaeger",
			Endpoint:  "",
			Collector: "",
			Service:   "store",
		},
		Datapath: path.Join(defaults.BaseDataPath(), "store"),
		Service: Service{
			Name:      "store",
			Namespace: "com.owncloud.api",
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
