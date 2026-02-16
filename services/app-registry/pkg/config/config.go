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
	Reva         *shared.Reva  `yaml:"reva"`

	AppRegistry AppRegistry `yaml:"app_registry"`

	Context context.Context `yaml:"-"`
}

type Log struct {
	Level  string `yaml:"level" env:"OCIS_LOG_LEVEL;APP_REGISTRY_LOG_LEVEL" desc:"The log level. Valid values are: 'panic', 'fatal', 'error', 'warn', 'info', 'debug', 'trace'." introductionVersion:"pre5.0"`
	Pretty bool   `yaml:"pretty" env:"OCIS_LOG_PRETTY;APP_REGISTRY_LOG_PRETTY" desc:"Activates pretty log output." introductionVersion:"pre5.0"`
	Color  bool   `yaml:"color" env:"OCIS_LOG_COLOR;APP_REGISTRY_LOG_COLOR" desc:"Activates colorized log output." introductionVersion:"pre5.0"`
	File   string `yaml:"file" env:"OCIS_LOG_FILE;APP_REGISTRY_LOG_FILE" desc:"The path to the log file. Activates logging to this file if set." introductionVersion:"pre5.0"`
}

type Service struct {
	Name string `yaml:"-"`
}

type Debug struct {
	Addr   string `yaml:"addr" env:"APP_REGISTRY_DEBUG_ADDR" desc:"Bind address of the debug server, where metrics, health, config and debug endpoints will be exposed." introductionVersion:"pre5.0"`
	Token  string `yaml:"token" env:"APP_REGISTRY_DEBUG_TOKEN" desc:"Token to secure the metrics endpoint." introductionVersion:"pre5.0"`
	Pprof  bool   `yaml:"pprof" env:"APP_REGISTRY_DEBUG_PPROF" desc:"Enables pprof, which can be used for profiling." introductionVersion:"pre5.0"`
	Zpages bool   `yaml:"zpages" env:"APP_REGISTRY_DEBUG_ZPAGES" desc:"Enables zpages, which can be used for collecting and viewing in-memory traces." introductionVersion:"pre5.0"`
}

type GRPCConfig struct {
	Addr      string                 `yaml:"addr" env:"APP_REGISTRY_GRPC_ADDR" desc:"The bind address of the GRPC service." introductionVersion:"pre5.0"`
	TLS       *shared.GRPCServiceTLS `yaml:"tls"`
	Namespace string                 `yaml:"-"`
	Protocol  string                 `yaml:"protocol" env:"OCIS_GRPC_PROTOCOL;APP_REGISTRY_GRPC_PROTOCOL" desc:"The transport protocol of the GRPC service." introductionVersion:"pre5.0"`
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
