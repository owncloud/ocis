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

	TokenManager *TokenManager `yaml:"token_manager"`
	Reva         *Reva         `yaml:"reva"`

	SkipUserGroupsInToken bool `yaml:"skip_user_groups_in_token" env:"OCDAV_SKIP_USER_GROUPS_IN_TOKEN"`

	WebdavNamespace string `yaml:"webdav_namespace" env:"OCDAV_WEBDAV_NAMESPACE"`
	FilesNamespace  string `yaml:"files_namespace" env:"OCDAV_FILES_NAMESPACE"`
	SharesNamespace string `yaml:"shares_namespace" env:"OCDAV_SHARES_NAMESPACE"`
	// PublicURL used to redirect /s/{token} URLs to
	PublicURL string `yaml:"public_url" env:"OCIS_URL;OCDAV_PUBLIC_URL"`

	// Insecure certificates allowed when making requests to the gateway
	Insecure bool `yaml:"insecure" env:"OCIS_INSECURE;OCDAV_INSECURE"`
	// Timeout in seconds when making requests to the gateway
	Timeout    int64      `yaml:"gateway_request_timeout" env:"OCDAV_GATEWAY_REQUEST_TIMEOUT"`
	Middleware Middleware `yaml:"middleware"`

	Context context.Context `yaml:"-"`
	Status  Status          `yaml:"-"`
}
type Tracing struct {
	Enabled   bool   `yaml:"enabled" env:"OCIS_TRACING_ENABLED;OCDAV_TRACING_ENABLED" desc:"Activates tracing."`
	Type      string `yaml:"type" env:"OCIS_TRACING_TYPE;OCDAV_TRACING_TYPE" desc:"The type of tracing. Defaults to \"\", which is the same as \"jaeger\". Allowed tracing types are \"jaeger\" and \"\" as of now."`
	Endpoint  string `yaml:"endpoint" env:"OCIS_TRACING_ENDPOINT;OCDAV_TRACING_ENDPOINT" desc:"The endpoint of the tracing agent."`
	Collector string `yaml:"collector" env:"OCIS_TRACING_COLLECTOR;OCDAV_TRACING_COLLECTOR" desc:"The HTTP endpoint for sending spans directly to a collector, i.e. http://jaeger-collector:14268/api/traces. Only used if the tracing endpoint is unset."`
}

type Log struct {
	Level  string `yaml:"level" env:"OCIS_LOG_LEVEL;OCDAV_LOG_LEVEL" desc:"The log level."`
	Pretty bool   `yaml:"pretty" env:"OCIS_LOG_PRETTY;OCDAV_LOG_PRETTY" desc:"Activates pretty log output."`
	Color  bool   `yaml:"color" env:"OCIS_LOG_COLOR;OCDAV_LOG_COLOR" desc:"Activates colorized log output."`
	File   string `yaml:"file" env:"OCIS_LOG_FILE;OCDAV_LOG_FILE" desc:"The target log file."`
}

type Service struct {
	Name string `yaml:"-"`
}

type Debug struct {
	Addr   string `yaml:"addr" env:"OCDAV_DEBUG_ADDR" desc:"Bind address of the debug server, where metrics, health, config and debug endpoints will be exposed."`
	Token  string `yaml:"token" env:"OCDAV_DEBUG_TOKEN" desc:"Token to secure the metrics endpoint"`
	Pprof  bool   `yaml:"pprof" env:"OCDAV_DEBUG_PPROF" desc:"Enables pprof, which can be used for profiling"`
	Zpages bool   `yaml:"zpages" env:"OCDAV_DEBUG_ZPAGES" desc:"Enables zpages, which can  be used for collecting and viewing traces in-memory."`
}

type HTTPConfig struct {
	Addr      string `yaml:"addr" env:"OCDAV_HTTP_ADDR" desc:"The address of the http service."`
	Namespace string `yaml:"-"`
	Protocol  string `yaml:"protocol" env:"OCDAV_HTTP_PROTOCOL" desc:"The transport protocol of the http service."`
	Prefix    string `yaml:"prefix" env:"OCDAV_HTTP_PREFIX"`
}

// Middleware configures reva middlewares.
type Middleware struct {
	Auth Auth `yaml:"auth"`
}

// Auth configures reva http auth middleware.
type Auth struct {
	CredentialsByUserAgent map[string]string `yaml:"credentials_by_user_agent"`
}

// Status holds the configurable values for the status.php
type Status struct {
	Version        string
	VersionString  string
	Product        string
	ProductName    string
	ProductVersion string
	Edition        string
}
