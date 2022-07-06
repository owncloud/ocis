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
	Type      string `yaml:"type" env:"OCIS_TRACING_TYPE;APP_REGISTRY_TRACING_TYPE" desc:"The type of tracing. Defaults to \"\", which is the same as \"jaeger\". Allowed tracing types are \"jaeger\" and \"\" as of now."`
	Endpoint  string `yaml:"endpoint" env:"OCIS_TRACING_ENDPOINT;APP_REGISTRY_TRACING_ENDPOINT"  desc:"The endpoint of the tracing agent."`
	Collector string `yaml:"collector" env:"OCIS_TRACING_COLLECTOR;APP_REGISTRY_TRACING_COLLECTOR" desc:"The HTTP endpoint for sending spans directly to a collector, i.e. http://jaeger-collector:14268/api/traces. Only used if the tracing endpoint is unset."`
}

type Log struct {
	Level  string `yaml:"level" env:"OCIS_LOG_LEVEL;APP_REGISTRY_LOG_LEVEL" desc:"The log level. Valid values are: \"panic\", \"fatal\", \"error\", \"warn\", \"info\", \"debug\", \"trace\"."`
	Pretty bool   `yaml:"pretty" env:"OCIS_LOG_PRETTY;APP_REGISTRY_LOG_PRETTY" desc:"Activates pretty log output."`
	Color  bool   `yaml:"color" env:"OCIS_LOG_COLOR;APP_REGISTRY_LOG_COLOR" desc:"Activates colorized log output."`
	File   string `yaml:"file" env:"OCIS_LOG_FILE;APP_REGISTRY_LOG_FILE" desc:"The path to the log file. Activates logging to this file if set."`
}

type Service struct {
	Name string `yaml:"-"`
}

type Debug struct {
	Addr   string `yaml:"addr" env:"APP_REGISTRY_DEBUG_ADDR" desc:"Bind address of the debug server, where metrics, health, config and debug endpoints will be exposed."`
	Token  string `yaml:"token" env:"APP_REGISTRY_DEBUG_TOKEN" desc:"Token to secure the metrics endpoint."`
	Pprof  bool   `yaml:"pprof" env:"APP_REGISTRY_DEBUG_PPROF" desc:"Enables pprof, which can be used for profiling."`
	Zpages bool   `yaml:"zpages" env:"APP_REGISTRY_DEBUG_ZPAGES" desc:"Enables zpages, which can be used for collecting and viewing in-memory traces."`
}

type GRPCConfig struct {
	Addr      string `yaml:"addr" env:"APP_REGISTRY_GRPC_ADDR" desc:"The bind address of the GRPC service."`
	Namespace string `yaml:"-"`
	Protocol  string `yaml:"protocol" env:"APP_REGISTRY_GRPC_PROTOCOL" desc:"The transport protocol of the GRPC service."`
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
