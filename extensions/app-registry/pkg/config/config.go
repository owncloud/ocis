package config

import (
	"context"

	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
)

type Config struct {
	Commons *shared.Commons `yaml:"-"` // don't use this directly as configuration for a service

	Service Service  `yaml:"-"`
	Tracing *Tracing `yaml:"tracing"`
	Log     *Log     `yaml:"log"`
	Debug   Debug    `yaml:"debug"`

	GRPC GRPCConfig `yaml:"grpc"`

	TokenManager *TokenManager `yaml:"token_manager"`
	Reva         *Reva         `yaml:"reva"`

	AppRegistry AppRegistry `yaml:"app_registry"`

	Supervised bool            `yaml:"-"`
	Context    context.Context `yaml:"-"`
}
type Tracing struct {
	Enabled   bool   `yaml:"enabled" env:"OCIS_TRACING_ENABLED;APP_REGISTRY_TRACING_ENABLED" desc:"Activates tracing."`
	Type      string `yaml:"type" env:"OCIS_TRACING_TYPE;APP_REGISTRY_TRACING_TYPE"`
	Endpoint  string `yaml:"endpoint" env:"OCIS_TRACING_ENDPOINT;APP_REGISTRY_TRACING_ENDPOINT" desc:"The endpoint to the tracing collector."`
	Collector string `yaml:"collector" env:"OCIS_TRACING_COLLECTOR;APP_REGISTRY_TRACING_COLLECTOR"`
}

type Log struct {
	Level  string `yaml:"level" env:"OCIS_LOG_LEVEL;APP_REGISTRY_LOG_LEVEL" desc:"The log level."`
	Pretty bool   `yaml:"pretty" env:"OCIS_LOG_PRETTY;APP_REGISTRY_LOG_PRETTY" desc:"Activates pretty log output."`
	Color  bool   `yaml:"color" env:"OCIS_LOG_COLOR;APP_REGISTRY_LOG_COLOR" desc:"Activates colorized log output."`
	File   string `yaml:"file" env:"OCIS_LOG_FILE;APP_REGISTRY_LOG_FILE" desc:"The target log file."`
}

type Service struct {
	Name string `yaml:"-"`
}

type Debug struct {
	Addr   string `yaml:"addr" env:"APP_REGISTRY_DEBUG_ADDR"`
	Token  string `yaml:"token" env:"APP_REGISTRY_DEBUG_TOKEN"`
	Pprof  bool   `yaml:"pprof" env:"APP_REGISTRY_DEBUG_PPROF"`
	Zpages bool   `yaml:"zpages" env:"APP_REGISTRY_DEBUG_ZPAGES"`
}

type GRPCConfig struct {
	Addr      string `yaml:"addr" env:"APP_REGISTRY_GRPC_ADDR" desc:"The address of the grpc service."`
	Namespace string `yaml:"-"`
	Protocol  string `yaml:"protocol" env:"APP_REGISTRY_GRPC_PROTOCOL" desc:"The transport protocol of the grpc service."`
}

type AppRegistry struct {
	MimeTypeConfig []MimeTypeConfig `yaml:"mimetypes"`
}

type MimeTypeConfig struct {
	MimeType      string `yaml:"mime_type" mapstructure:"mime_type"`
	Extension     string `yaml:"extension" mapstructure:"extension"`
	Name          string `yaml:"name" mapstructure:"name"`
	Description   string `yaml:"description" mapstructure:"description"`
	Icon          string `yaml:"icon" mapstructure:"icon"`
	DefaultApp    string `yaml:"default_app" mapstructure:"default_app"`
	AllowCreation bool   `yaml:"allow_creation" mapstructure:"allow_creation"`
}
