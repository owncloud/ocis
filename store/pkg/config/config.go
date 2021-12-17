package config

import (
	"context"
	"path"

	"github.com/owncloud/ocis/ocis-pkg/config/defaults"
)

// Debug defines the available debug configuration.
type Debug struct {
	Addr   string `ocisConfig:"addr" env:"STORE_DEBUG_ADDR"`
	Token  string `ocisConfig:"token" env:"STORE_DEBUG_TOKEN"`
	Pprof  bool   `ocisConfig:"pprof" env:"STORE_DEBUG_PPROF"`
	Zpages bool   `ocisConfig:"zpages" env:"STORE_DEBUG_ZPAGES"`
}

// GRPC defines the available grpc configuration.
type GRPC struct {
	Addr      string `ocisConfig:"addr" env:"STORE_GRPC_ADDR"`
	Namespace string
}

// Service defines the available service configuration.
type Service struct {
	Name    string
	Version string
}

// Tracing defines the available tracing configuration.
type Tracing struct {
	Enabled   bool   `ocisConfig:"enabled" env:"OCIS_TRACING_ENABLED;STORE_TRACING_ENABLED"`
	Type      string `ocisConfig:"type" env:"OCIS_TRACING_TYPE;STORE_TRACING_TYPE"`
	Endpoint  string `ocisConfig:"endpoint" env:"OCIS_TRACING_ENDPOINT;STORE_TRACING_ENDPOINT"`
	Collector string `ocisConfig:"collector" env:"OCIS_TRACING_COLLECTOR;STORE_TRACING_COLLECTOR"`
	Service   string `ocisConfig:"service" env:"STORE_TRACING_SERVICE"` //TODO: should this be an ID? or the same as Service.Name?
}

// Log defines the available log configuration.
type Log struct {
	Level  string `mapstructure:"level" env:"OCIS_LOG_LEVEL;STORE_LOG_LEVEL"`
	Pretty bool   `mapstructure:"pretty" env:"OCIS_LOG_PRETTY;STORE_LOG_PRETTY"`
	Color  bool   `mapstructure:"color" env:"OCIS_LOG_COLOR;STORE_LOG_COLOR"`
	File   string `mapstructure:"file" env:"OCIS_LOG_FILE;STORE_LOG_FILE"`
}

// Config combines all available configuration parts.
type Config struct {
	Service Service `ocisConfig:"service"`

	Tracing Tracing `ocisConfig:"tracing"`
	Log     Log     `ocisConfig:"log"`
	Debug   Debug   `ocisConfig:"debug"`

	GRPC GRPC `ocisConfig:"grpc"`

	Datapath string `ocisConfig:"data_path" env:"STORE_DATA_PATH"`

	Context    context.Context
	Supervised bool
}

func DefaultConfig() *Config {
	return &Config{
		Debug: Debug{
			Addr:   "127.0.0.1:9464",
			Token:  "",
			Pprof:  false,
			Zpages: false,
		},
		GRPC: GRPC{
			Addr:      "127.0.0.1:9460",
			Namespace: "com.owncloud.api",
		},
		Service: Service{
			Name: "store",
		},
		Tracing: Tracing{
			Enabled:   false,
			Type:      "jaeger",
			Endpoint:  "",
			Collector: "",
			Service:   "store",
		},
		Datapath: path.Join(defaults.BaseDataPath(), "store"),
	}
}
