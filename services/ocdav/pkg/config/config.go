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
	Reva         *shared.Reva  `yaml:"reva"`

	SkipUserGroupsInToken bool `yaml:"skip_user_groups_in_token" env:"OCDAV_SKIP_USER_GROUPS_IN_TOKEN" desc:"Disables the loading of user's group memberships from the reva access token." introductionVersion:"pre5.0"`

	WebdavNamespace string `yaml:"webdav_namespace" env:"OCDAV_WEBDAV_NAMESPACE" desc:"Jail requests to /dav/webdav into this CS3 namespace. Supports template layouting with CS3 User properties." introductionVersion:"pre5.0"`
	FilesNamespace  string `yaml:"files_namespace" env:"OCDAV_FILES_NAMESPACE" desc:"Jail requests to /dav/files/{username} into this CS3 namespace. Supports template layouting with CS3 User properties." introductionVersion:"pre5.0"`
	SharesNamespace string `yaml:"shares_namespace" env:"OCDAV_SHARES_NAMESPACE" desc:"The human readable path for the share jail. Relative to a users personal space root. Upcased intentionally." introductionVersion:"pre5.0"`
	OCMNamespace    string `yaml:"ocm_namespace" env:"OCDAV_OCM_NAMESPACE" desc:"The human readable path prefix for the ocm shares." introductionVersion:"5.0"`
	// PublicURL used to redirect /s/{token} URLs to
	PublicURL string `yaml:"public_url" env:"OCIS_URL;OCDAV_PUBLIC_URL" desc:"URL where oCIS is reachable for users." introductionVersion:"pre5.0"`

	// Insecure certificates allowed when making requests to the gateway
	Insecure bool `yaml:"insecure" env:"OCIS_INSECURE;OCDAV_INSECURE" desc:"Allow insecure connections to the GATEWAY service." introductionVersion:"pre5.0"`
	// Timeout in seconds when making requests to the gateway
	Timeout int64 `yaml:"gateway_request_timeout" env:"OCDAV_GATEWAY_REQUEST_TIMEOUT" desc:"Request timeout in seconds for requests from the oCDAV service to the GATEWAY service." introductionVersion:"pre5.0"`

	MachineAuthAPIKey string `yaml:"machine_auth_api_key" env:"OCIS_MACHINE_AUTH_API_KEY;OCDAV_MACHINE_AUTH_API_KEY" desc:"Machine auth API key used to validate internal requests necessary for the access to resources from other services." introductionVersion:"pre5.0"`

	Context context.Context `yaml:"-"`
	Status  Status          `yaml:"-"`

	AllowPropfindDepthInfinity bool `yaml:"allow_propfind_depth_infinity" env:"OCDAV_ALLOW_PROPFIND_DEPTH_INFINITY" desc:"Allow the use of depth infinity in PROPFINDS. When enabled, a propfind will traverse through all subfolders. If many subfolders are expected, depth infinity can cause heavy server load and/or delayed response times." introductionVersion:"pre5.0"`
}

type Log struct {
	Level  string `yaml:"level" env:"OCIS_LOG_LEVEL;OCDAV_LOG_LEVEL" desc:"The log level. Valid values are: 'panic', 'fatal', 'error', 'warn', 'info', 'debug', 'trace'." introductionVersion:"pre5.0"`
	Pretty bool   `yaml:"pretty" env:"OCIS_LOG_PRETTY;OCDAV_LOG_PRETTY" desc:"Activates pretty log output." introductionVersion:"pre5.0"`
	Color  bool   `yaml:"color" env:"OCIS_LOG_COLOR;OCDAV_LOG_COLOR" desc:"Activates colorized log output." introductionVersion:"pre5.0"`
	File   string `yaml:"file" env:"OCIS_LOG_FILE;OCDAV_LOG_FILE" desc:"The path to the log file. Activates logging to this file if set." introductionVersion:"pre5.0"`
}

type Service struct {
	Name string `yaml:"-"`
}

type Debug struct {
	Addr   string `yaml:"addr" env:"OCDAV_DEBUG_ADDR" desc:"Bind address of the debug server, where metrics, health, config and debug endpoints will be exposed." introductionVersion:"pre5.0"`
	Token  string `yaml:"token" env:"OCDAV_DEBUG_TOKEN" desc:"Token to secure the metrics endpoint." introductionVersion:"pre5.0"`
	Pprof  bool   `yaml:"pprof" env:"OCDAV_DEBUG_PPROF" desc:"Enables pprof, which can be used for profiling." introductionVersion:"pre5.0"`
	Zpages bool   `yaml:"zpages" env:"OCDAV_DEBUG_ZPAGES" desc:"Enables zpages, which can be used for collecting and viewing in-memory traces." introductionVersion:"pre5.0"`
}

type HTTPConfig struct {
	Addr      string `yaml:"addr" env:"OCDAV_HTTP_ADDR" desc:"The bind address of the HTTP service." introductionVersion:"pre5.0"`
	Namespace string `yaml:"-"`
	Protocol  string `yaml:"protocol" env:"OCDAV_HTTP_PROTOCOL" desc:"The transport protocol of the HTTP service." introductionVersion:"pre5.0"`
	Prefix    string `yaml:"prefix" env:"OCDAV_HTTP_PREFIX" desc:"A URL path prefix for the handler." introductionVersion:"pre5.0"`
	CORS      `yaml:"cors"`
}

// CORS defines the available cors configuration.
type CORS struct {
	AllowedOrigins   []string `yaml:"allow_origins" env:"OCIS_CORS_ALLOW_ORIGINS;OCDAV_CORS_ALLOW_ORIGINS" desc:"A list of allowed CORS origins. See following chapter for more details: *Access-Control-Allow-Origin* at https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Origin. See the Environment Variable Types description for more details." introductionVersion:"pre5.0"`
	AllowedMethods   []string `yaml:"allow_methods" env:"OCIS_CORS_ALLOW_METHODS;OCDAV_CORS_ALLOW_METHODS" desc:"A list of allowed CORS methods. See following chapter for more details: *Access-Control-Request-Method* at https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Request-Method. See the Environment Variable Types description for more details." introductionVersion:"pre5.0"`
	AllowedHeaders   []string `yaml:"allow_headers" env:"OCIS_CORS_ALLOW_HEADERS;OCDAV_CORS_ALLOW_HEADERS" desc:"A list of allowed CORS headers. See following chapter for more details: *Access-Control-Request-Headers* at https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Request-Headers. See the Environment Variable Types description for more details." introductionVersion:"pre5.0"`
	AllowCredentials bool     `yaml:"allow_credentials" env:"OCIS_CORS_ALLOW_CREDENTIALS;OCDAV_CORS_ALLOW_CREDENTIALS" desc:"Allow credentials for CORS.See following chapter for more details: *Access-Control-Allow-Credentials* at https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Credentials." introductionVersion:"pre5.0"`
}

// Status holds the configurable values for the status.php
type Status struct {
	Version        string
	VersionString  string
	Product        string
	ProductName    string
	ProductVersion string
	Edition        string `yaml:"edition" env:"OCIS_EDITION;OCDAV_EDITION" desc:"Edition of oCIS. Used for branding purposes." introductionVersion:"pre5.0"`
}
