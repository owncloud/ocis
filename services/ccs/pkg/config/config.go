package config

import (
	"context"

	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
)

type Config struct {
	Commons *shared.Commons `yaml:"-"` // don't use this directly as configuration for a service
	Service Service         `yaml:"-"`
	Tracing *Tracing        `yaml:"tracing"`
	Log     *Log            `yaml:"log"`
	Debug   Debug           `yaml:"debug"`

	HTTP HTTPConfig `yaml:"http"`

	Context   context.Context `yaml:"-"`
	JWTSecret string          `yaml:"jwt_secret" env:"OCIS_JWT_SECRET;CALDAV_JWT_SECRET" desc:"The secret to mint and validate jwt tokens."`
}

type Log struct {
	Level  string `yaml:"level" env:"OCIS_LOG_LEVEL;CALDAV_LOG_LEVEL" desc:"The log level. Valid values are: 'panic', 'fatal', 'error', 'warn', 'info', 'debug', 'trace'."`
	Pretty bool   `yaml:"pretty" env:"OCIS_LOG_PRETTY;CALDAV_LOG_PRETTY" desc:"Activates pretty log output."`
	Color  bool   `yaml:"color" env:"OCIS_LOG_COLOR;CALDAV_LOG_COLOR" desc:"Activates colorized log output."`
	File   string `yaml:"file" env:"OCIS_LOG_FILE;CALDAV_LOG_FILE" desc:"The path to the log file. Activates logging to this file if set."`
}

type Service struct {
	Name string `yaml:"-"`
}

type Debug struct {
	Addr   string `yaml:"addr" env:"CALDAV_DEBUG_ADDR" desc:"Bind address of the debug server, where metrics, health, config and debug endpoints will be exposed."`
	Token  string `yaml:"token" env:"CALDAV_DEBUG_TOKEN" desc:"Token to secure the metrics endpoint."`
	Pprof  bool   `yaml:"pprof" env:"CALDAV_DEBUG_PPROF" desc:"Enables pprof, which can be used for profiling."`
	Zpages bool   `yaml:"zpages" env:"CALDAV_DEBUG_ZPAGES" desc:"Enables zpages, which can be used for collecting and viewing in-memory traces."`
}

type HTTPConfig struct {
	Addr      string                `yaml:"addr" env:"CALDAV_HTTP_ADDR" desc:"The bind address of the HTTP service."`
	Namespace string                `yaml:"-"`
	Root      string                `yaml:"root" env:"CALDAV_HTTP_ROOT" desc:"Subdirectory that serves as the root for this HTTP service."`
	TLS       shared.HTTPServiceTLS `yaml:"tls"`
	Protocol  string                `yaml:"protocol" env:"CALDAV_HTTP_PROTOCOL" desc:"The transport protocol of the HTTP service."`
	Prefix    string                `yaml:"prefix" env:"CALDAV_HTTP_PREFIX" desc:"A URL path prefix for the handler."`
	CORS      `yaml:"cors"`
}

// CORS defines the available cors configuration.
type CORS struct {
	AllowedOrigins   []string `yaml:"allow_origins" env:"OCIS_CORS_ALLOW_ORIGINS;CALDAV_CORS_ALLOW_ORIGINS" desc:"A list of allowed CORS origins. See following chapter for more details: *Access-Control-Allow-Origin* at https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Origin. See the Environment Variable Types description for more details."`
	AllowedMethods   []string `yaml:"allow_methods" env:"OCIS_CORS_ALLOW_METHODS;CALDAV_CORS_ALLOW_METHODS" desc:"A list of allowed CORS methods. See following chapter for more details: *Access-Control-Request-Method* at https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Request-Method. See the Environment Variable Types description for more details."`
	AllowedHeaders   []string `yaml:"allow_headers" env:"OCIS_CORS_ALLOW_HEADERS;CALDAV_CORS_ALLOW_HEADERS" desc:"A list of allowed CORS headers. See following chapter for more details: *Access-Control-Request-Headers* at https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Request-Headers. See the Environment Variable Types description for more details."`
	AllowCredentials bool     `yaml:"allow_credentials" env:"OCIS_CORS_ALLOW_CREDENTIALS;CALDAV_CORS_ALLOW_CREDENTIALS" desc:"Allow credentials for CORS.See following chapter for more details: *Access-Control-Allow-Credentials* at https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Credentials."`
}

// Status holds the configurable values for the status.php
type Status struct {
	Version        string
	VersionString  string
	Product        string
	ProductName    string
	ProductVersion string
	Edition        string `yaml:"edition" env:"OCIS_EDITION;CALDAV_EDITION"`
}
