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

	// JWTSecret used to verify reva access token

	TransferSecret string `yaml:"transfer_secret" env:"STORAGE_TRANSFER_SECRET" desc:"Transfer secret for signing file up- and download requests."`

	TokenManager      *TokenManager `yaml:"token_manager"`
	Reva              *Reva         `yaml:"reva"`
	MachineAuthAPIKey string        `yaml:"machine_auth_api_key" env:"OCIS_MACHINE_AUTH_API_KEY;FRONTEND_MACHINE_AUTH_API_KEY" desc:"Machine auth API key used to validate internal requests necessary to access resources from other services."`

	SkipUserGroupsInToken bool `yaml:"skip_user_groups_in_token" env:"FRONTEND_SKIP_USER_GROUPS_IN_TOKEN" desc:"Disables the loading of user's group memberships from the reva access token."`

	EnableFavorites          bool   `yaml:"enable_favorites" env:"FRONTEND_ENABLE_FAVORITES" desc:"Enables the support for favorites in the frontend."`
	EnableProjectSpaces      bool   `yaml:"enable_project_spaces" env:"FRONTEND_ENABLE_PROJECT_SPACES" desc:"Indicates to clients that project spaces are supposed to be made available."`
	EnableShareJail          bool   `yaml:"enable_share_jail" env:"FRONTEND_ENABLE_SHARE_JAIL" desc:"Indicates to clients that the share jail is supposed to be used."`
	UploadMaxChunkSize       int    `yaml:"upload_max_chunk_size" env:"FRONTEND_UPLOAD_MAX_CHUNK_SIZE" desc:"Sets the max chunk sizes for uploads via the frontend." `
	UploadHTTPMethodOverride string `yaml:"upload_http_method_override" env:"FRONTEND_UPLOAD_HTTP_METHOD_OVERRIDE" desc:"Advise TUS to replace PATCH requests by POST requests."`
	DefaultUploadProtocol    string `yaml:"default_upload_protocol" env:"FRONTEND_DEFAULT_UPLOAD_PROTOCOL" desc:"The default upload protocol to use in the frontend (e.g. tus)."`
	EnableResharing          bool   `yaml:"enable_resharing" env:"FRONTEND_ENABLE_RESHARING" desc:"Enables the support for resharing in the frontend."`

	PublicURL string `yaml:"public_url" env:"OCIS_URL;FRONTEND_PUBLIC_URL" desc:"The public facing URL of the oCIS frontend."`

	AppHandler  AppHandler  `yaml:"app_handler"`
	Archiver    Archiver    `yaml:"archiver"`
	DataGateway DataGateway `yaml:"data_gateway"`
	OCS         OCS         `yaml:"ocs"`
	Checksums   Checksums   `yaml:"checksums"`

	Middleware Middleware `yaml:"middleware"`

	Supervised bool            `yaml:"-"`
	Context    context.Context `yaml:"-"`
}
type Tracing struct {
	Enabled   bool   `yaml:"enabled" env:"OCIS_TRACING_ENABLED;FRONTEND_TRACING_ENABLED" desc:"Activates tracing."`
	Type      string `yaml:"type" env:"OCIS_TRACING_TYPE;FRONTEND_TRACING_TYPE" desc:"The type of tracing. Defaults to \"\", which is the same as \"jaeger\". Allowed tracing types are \"jaeger\" and \"\" as of now."`
	Endpoint  string `yaml:"endpoint" env:"OCIS_TRACING_ENDPOINT;FRONTEND_TRACING_ENDPOINT" desc:"The endpoint of the tracing agent."`
	Collector string `yaml:"collector" env:"OCIS_TRACING_COLLECTOR;FRONTEND_TRACING_COLLECTOR" desc:"The HTTP endpoint for sending spans directly to a collector, i.e. http://jaeger-collector:14268/api/traces. Only used if the tracing endpoint is unset."`
}

type Log struct {
	Level  string `yaml:"level" env:"OCIS_LOG_LEVEL;FRONTEND_LOG_LEVEL" desc:"The log level. Valid values are: \"panic\", \"fatal\", \"error\", \"warn\", \"info\", \"debug\", \"trace\"."`
	Pretty bool   `yaml:"pretty" env:"OCIS_LOG_PRETTY;FRONTEND_LOG_PRETTY" desc:"Activates pretty log output."`
	Color  bool   `yaml:"color" env:"OCIS_LOG_COLOR;FRONTEND_LOG_COLOR" desc:"Activates colorized log output."`
	File   string `yaml:"file" env:"OCIS_LOG_FILE;FRONTEND_LOG_FILE" desc:"The path to the log file. Activates logging to this file if set."`
}

type Service struct {
	Name string `yaml:"-"`
}

type Debug struct {
	Addr   string `yaml:"addr" env:"FRONTEND_DEBUG_ADDR" desc:"Bind address of the debug server, where metrics, health, config and debug endpoints will be exposed."`
	Token  string `yaml:"token" env:"FRONTEND_DEBUG_TOKEN" desc:"Token to secure the metrics endpoint."`
	Pprof  bool   `yaml:"pprof" env:"FRONTEND_DEBUG_PPROF" desc:"Enables pprof, which can be used for profiling."`
	Zpages bool   `yaml:"zpages" env:"FRONTEND_DEBUG_ZPAGES" desc:"Enables zpages, which can be used for collecting and viewing in-memory traces."`
}

type HTTPConfig struct {
	Addr      string `yaml:"addr" env:"FRONTEND_HTTP_ADDR" desc:"The bind address of the HTTP service."`
	Namespace string `yaml:"-"`
	Protocol  string `yaml:"protocol" env:"FRONTEND_HTTP_PROTOCOL" desc:"The transport protocol of the HTTP service."`
	Prefix    string `yaml:"prefix" env:"FRONTEND_HTTP_PREFIX" desc:"The Path prefix where the frontend can be accessed (defaults to /)."`
}

// Middleware configures reva middlewares.
type Middleware struct {
	Auth Auth `yaml:"auth"`
}

// Auth configures reva http auth middleware.
type Auth struct {
	CredentialsByUserAgent map[string]string `yaml:"credentials_by_user_agent"`
}

type AppHandler struct {
	Prefix   string `yaml:"-"`
	Insecure bool   `yaml:"insecure" env:"OCIS_INSECURE;FRONTEND_APP_HANDLER_INSECURE" desc:"Allow insecure connections to the frontend."`
}

type Archiver struct {
	MaxNumFiles int64  `yaml:"max_num_files" env:"FRONTEND_ARCHIVER_MAX_NUM_FILES" desc:"Max number of files that can be packed into an archive."`
	MaxSize     int64  `yaml:"max_size" env:"FRONTEND_ARCHIVER_MAX_SIZE" desc:"Max size of the zip archive the archiver can create."`
	Prefix      string `yaml:"-"`
	Insecure    bool   `yaml:"insecure" env:"OCIS_INSECURE;FRONTEND_ARCHIVER_INSECURE" desc:"Allow insecure connections to the archiver."`
}

type DataGateway struct {
	Prefix string `yaml:"prefix" env:"FRONTEND_DATA_GATEWAY_PREFIX"`
}

type OCS struct {
	Prefix                  string             `yaml:"prefix" env:"FRONTEND_OCS_PREFIX" desc:"Path prefix for the OCS service"`
	SharePrefix             string             `yaml:"share_prefix" env:"FRONTEND_OCS_SHARE_PREFIX" desc:"Path prefix for shares."`
	HomeNamespace           string             `yaml:"home_namespace" env:"FRONTEND_OCS_HOME_NAMESPACE" desc:"Homespace namespace identifier."`
	AdditionalInfoAttribute string             `yaml:"additional_info_attribute" env:"FRONTEND_OCS_ADDITIONAL_INFO_ATTRIBUTE" desc:"Additional information attribute for the user like {{.Mail}}."`
	ResourceInfoCacheTTL    int                `yaml:"resource_info_cache_ttl" env:"FRONTEND_OCS_RESOURCE_INFO_CACHE_TTL" desc:"Max TTL for the resource info cache."`
	CacheWarmupDriver       string             `yaml:"cache_warmup_driver,omitempty"`  // not supported by the oCIS product, therefore not part of docs
	CacheWarmupDrivers      CacheWarmupDrivers `yaml:"cache_warmup_drivers,omitempty"` // not supported by the oCIS product, therefore not part of docs
}

type CacheWarmupDrivers struct {
	CBOX CBOXDriver `yaml:"cbox,omitempty"`
}

type CBOXDriver struct {
	DBUsername string `yaml:"db_username,omitempty"`
	DBPassword string `yaml:"db_password,omitempty"`
	DBHost     string `yaml:"db_host,omitempty"`
	DBPort     int    `yaml:"db_port,omitempty"`
	DBName     string `yaml:"db_name,omitempty"`
	Namespace  string `yaml:"namespace,omitempty"`
}

type Checksums struct {
	SupportedTypes      []string `yaml:"supported_types" env:"FRONTEND_CHECKSUMS_SUPPORTED_TYPES" desc:"Supported checksum types to be announced to the client. You can provide multiple types separated by blank or comma."`
	PreferredUploadType string   `yaml:"preferred_upload_type" env:"FRONTEND_CHECKSUMS_PREFERRED_UPLOAD_TYPES" desc:"Preferred checksum types to be announced to the client for uploads (e.g. md5)"`
}
