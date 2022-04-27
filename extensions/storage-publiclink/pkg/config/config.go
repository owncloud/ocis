package config

import (
	"context"

	"github.com/owncloud/ocis/ocis-pkg/shared"
)

type Config struct {
	*shared.Commons `yaml:"-"`
	Service         Service  `yaml:"-"`
	Tracing         *Tracing `yaml:"tracing,omitempty"`
	Logging         *Logging `yaml:"log,omitempty"`
	Debug           Debug    `yaml:"debug,omitempty"`
	Supervised      bool     `yaml:"supervised,omitempty"`

	GRPC GRPCConfig `yaml:"grpc,omitempty"`

	Context               context.Context `yaml:"context,omitempty"`
	JWTSecret             string          `yaml:"jwt_secret,omitempty"`
	GatewayEndpoint       string          `yaml:"gateway_endpoint,omitempty"`
	SkipUserGroupsInToken bool            `yaml:"skip_user_groups_in_token,omitempty"`
	AuthProvider          AuthProvider    `yaml:"auth_provider,omitempty"`
	StorageProvider       StorageProvider `yaml:"storage_provider,omitempty"`
}
type Tracing struct {
	Enabled   bool   `yaml:"enabled" env:"OCIS_TRACING_ENABLED;STORAGE_METADATA_TRACING_ENABLED" desc:"Activates tracing."`
	Type      string `yaml:"type" env:"OCIS_TRACING_TYPE;STORAGE_METADATA_TRACING_TYPE"`
	Endpoint  string `yaml:"endpoint" env:"OCIS_TRACING_ENDPOINT;STORAGE_METADATA_TRACING_ENDPOINT" desc:"The endpoint to the tracing collector."`
	Collector string `yaml:"collector" env:"OCIS_TRACING_COLLECTOR;STORAGE_METADATA_TRACING_COLLECTOR"`
}

type Logging struct {
	Level  string `yaml:"level" env:"OCIS_LOG_LEVEL;STORAGE_METADATA_LOG_LEVEL" desc:"The log level."`
	Pretty bool   `yaml:"pretty" env:"OCIS_LOG_PRETTY;STORAGE_METADATA_LOG_PRETTY" desc:"Activates pretty log output."`
	Color  bool   `yaml:"color" env:"OCIS_LOG_COLOR;STORAGE_METADATA_LOG_COLOR" desc:"Activates colorized log output."`
	File   string `yaml:"file" env:"OCIS_LOG_FILE;STORAGE_METADATA_LOG_FILE" desc:"The target log file."`
}

type Service struct {
	Name string `yaml:"-"`
}

type Debug struct {
	Addr   string `yaml:"addr" env:"STORAGE_METADATA_DEBUG_ADDR"`
	Token  string `yaml:"token" env:"STORAGE_METADATA_DEBUG_TOKEN"`
	Pprof  bool   `yaml:"pprof" env:"STORAGE_METADATA_DEBUG_PPROF"`
	Zpages bool   `yaml:"zpages" env:"STORAGE_METADATA_DEBUG_ZPAGES"`
}

type GRPCConfig struct {
	Addr     string `yaml:"addr" env:"STORAGE_METADATA_GRPC_ADDR" desc:"The address of the grpc service."`
	Protocol string `yaml:"protocol" env:"STORAGE_METADATA_GRPC_PROTOCOL" desc:"The transport protocol of the grpc service."`
}

type AuthProvider struct {
	GatewayEndpoint string
}

type StorageProvider struct {
	MountID         string
	GatewayEndpoint string
}
