package config

import (
	"context"

	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
)

// Config defines the root config structure
type Config struct {
	Commons *shared.Commons `yaml:"-"` // don't use this directly as configuration for a service
	Service Service         `yaml:"-"`
	Tracing *Tracing        `yaml:"tracing"`
	Log     *Log            `yaml:"log"`
	Debug   Debug           `yaml:"debug"`

	GRPC GRPCConfig `yaml:"grpc"`
	HTTP HTTP       `yaml:"http"`

	GRPCClientTLS *shared.GRPCClientTLS `yaml:"grpc_client_tls"`

	TokenManager *TokenManager `yaml:"token_manager"`
	Reva         *shared.Reva  `yaml:"reva"`

	SkipUserGroupsInToken bool `yaml:"skip_user_groups_in_token" env:"AUTH_APP_SKIP_USER_GROUPS_IN_TOKEN" desc:"Disables the encoding of the user's group memberships in the access token. This reduces the token size, especially when users are members of a large number of groups." introductionVersion:"7.0.0"`

	MachineAuthAPIKey string `yaml:"machine_auth_api_key" env:"OCIS_MACHINE_AUTH_API_KEY;AUTH_APP_MACHINE_AUTH_API_KEY" desc:"The machine auth API key used to validate internal requests necessary to access resources from other services." introductionVersion:"7.0.0"`

	AllowImpersonation bool `yaml:"allow_impersonation" env:"AUTH_APP_ENABLE_IMPERSONATION" desc:"Allows admins to create app tokens for other users. Used for migration. Do NOT use in productive deployments." introductionVersion:"7.0.0"`

	Supervised bool            `yaml:"-"`
	Context    context.Context `yaml:"-"`
}

// Log defines the loging configuration
type Log struct {
	Level  string `yaml:"level" env:"OCIS_LOG_LEVEL;AUTH_APP_LOG_LEVEL" desc:"The log level. Valid values are: 'panic', 'fatal', 'error', 'warn', 'info', 'debug', 'trace'." introductionVersion:"7.0.0"`
	Pretty bool   `yaml:"pretty" env:"OCIS_LOG_PRETTY;AUTH_APP_LOG_PRETTY" desc:"Activates pretty log output." introductionVersion:"7.0.0"`
	Color  bool   `yaml:"color" env:"OCIS_LOG_COLOR;AUTH_APP_LOG_COLOR" desc:"Activates colorized log output." introductionVersion:"7.0.0"`
	File   string `yaml:"file" env:"OCIS_LOG_FILE;AUTH_APP_LOG_FILE" desc:"The path to the log file. Activates logging to this file if set." introductionVersion:"7.0.0"`
}

// Service defines the service configuration
type Service struct {
	Name string `yaml:"-"`
}

// Debug defines the debug configuration
type Debug struct {
	Addr   string `yaml:"addr" env:"AUTH_APP_DEBUG_ADDR" desc:"Bind address of the debug server, where metrics, health, config and debug endpoints will be exposed." introductionVersion:"7.0.0"`
	Token  string `yaml:"token" env:"AUTH_APP_DEBUG_TOKEN" desc:"Token to secure the metrics endpoint." introductionVersion:"7.0.0"`
	Pprof  bool   `yaml:"pprof" env:"AUTH_APP_DEBUG_PPROF" desc:"Enables pprof, which can be used for profiling." introductionVersion:"7.0.0"`
	Zpages bool   `yaml:"zpages" env:"AUTH_APP_DEBUG_ZPAGES" desc:"Enables zpages, which can  be used for collecting and viewing traces in-memory." introductionVersion:"7.0.0"`
}

// GRPCConfig defines the GRPC configuration
type GRPCConfig struct {
	Addr      string                 `yaml:"addr" env:"AUTH_APP_GRPC_ADDR" desc:"The bind address of the GRPC service." introductionVersion:"7.0.0"`
	TLS       *shared.GRPCServiceTLS `yaml:"tls"`
	Namespace string                 `yaml:"-"`
	Protocol  string                 `yaml:"protocol" env:"OCIS_GRPC_PROTOCOL;AUTH_APP_GRPC_PROTOCOL" desc:"The transport protocol of the GRPC service." introductionVersion:"7.0.0"`
}

// HTTP defines the available http configuration.
type HTTP struct {
	Addr      string                `yaml:"addr" env:"AUTH_APP_HTTP_ADDR" desc:"The bind address of the HTTP service." introductionVersion:"pre5.0"`
	Namespace string                `yaml:"-"`
	Root      string                `yaml:"root" env:"AUTH_APP_HTTP_ROOT" desc:"Subdirectory that serves as the root for this HTTP service." introductionVersion:"pre5.0"`
	CORS      CORS                  `yaml:"cors"`
	TLS       shared.HTTPServiceTLS `yaml:"tls"`
}

// CORS defines the available cors configuration.
type CORS struct {
	AllowedOrigins   []string `yaml:"allow_origins" env:"OCIS_CORS_ALLOW_ORIGINS;AUTH_APP_CORS_ALLOW_ORIGINS" desc:"A list of allowed CORS origins. See following chapter for more details: *Access-Control-Allow-Origin* at https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Origin. See the Environment Variable Types description for more details." introductionVersion:"pre5.0"`
	AllowedMethods   []string `yaml:"allow_methods" env:"OCIS_CORS_ALLOW_METHODS;AUTH_APP_CORS_ALLOW_METHODS" desc:"A list of allowed CORS methods. See following chapter for more details: *Access-Control-Request-Method* at https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Request-Method. See the Environment Variable Types description for more details." introductionVersion:"pre5.0"`
	AllowedHeaders   []string `yaml:"allow_headers" env:"OCIS_CORS_ALLOW_HEADERS;AUTH_APP_CORS_ALLOW_HEADERS" desc:"A list of allowed CORS headers. See following chapter for more details: *Access-Control-Request-Headers* at https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Request-Headers. See the Environment Variable Types description for more details." introductionVersion:"pre5.0"`
	AllowCredentials bool     `yaml:"allow_credentials" env:"OCIS_CORS_ALLOW_CREDENTIALS;AUTH_APP_CORS_ALLOW_CREDENTIALS" desc:"Allow credentials for CORS.See following chapter for more details: *Access-Control-Allow-Credentials* at https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Credentials." introductionVersion:"pre5.0"`
}
