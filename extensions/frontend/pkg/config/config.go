package config

import "github.com/owncloud/ocis/ocis-pkg/shared"

type Config struct {
	*shared.Commons `yaml:"-"`
	Service         Service  `yaml:"-"`
	Tracing         *Tracing `yaml:"tracing,omitempty"`
	Logging         *Logging `yaml:"log,omitempty"`
	Debug           Debug    `yaml:"debug,omitempty"`
	Supervised      bool     `yaml:"-"`

	HTTP HTTPConfig `yaml:"http,omitempty"`

	// JWTSecret used to verify reva access token

	TransferSecret string `yaml:"transfer_secret,omitempty" env:"STORAGE_TRANSFER_SECRET"`

	TokenManager *TokenManager `yaml:"token_manager,omitempty"`
	Reva         *Reva         `yaml:"reva,omitempty"`

	SkipUserGroupsInToken bool `yaml:"skip_users_groups_in_token,omitempty"`

	EnableFavorites          bool   `yaml:"favorites,omitempty"`
	EnableProjectSpaces      bool   `yaml:"enable_project_spaces,omitempty"`
	UploadMaxChunkSize       int    `yaml:"upload_max_chunk_size,omitempty"`
	UploadHTTPMethodOverride string `yaml:"upload_http_method_override,omitempty"`
	DefaultUploadProtocol    string `yaml:"default_upload_protocol,omitempty"`

	PublicURL string `yaml:"public_url,omitempty" env:"OCIS_URL;FRONTEND_PUBLIC_URL"`

	Archiver    Archiver    `yaml:"archiver,omitempty"`
	AppProvider AppProvider `yaml:"app_provider,omitempty"`
	DataGateway DataGateway `yaml:"data_gateway,omitempty"`
	OCS         OCS         `yaml:"ocs,omitempty"`
	AuthMachine AuthMachine `yaml:"auth_machine,omitempty"`
	Checksums   Checksums   `yaml:"checksums,omitempty"`

	Middleware Middleware `yaml:"middleware,omitempty"`
}
type Tracing struct {
	Enabled   bool   `yaml:"enabled,omitempty" env:"OCIS_TRACING_ENABLED;FRONTEND_TRACING_ENABLED" desc:"Activates tracing."`
	Type      string `yaml:"type,omitempty" env:"OCIS_TRACING_TYPE;FRONTEND_TRACING_TYPE"`
	Endpoint  string `yaml:"endpoint,omitempty" env:"OCIS_TRACING_ENDPOINT;FRONTEND_TRACING_ENDPOINT" desc:"The endpoint to the tracing collector."`
	Collector string `yaml:"collector,omitempty" env:"OCIS_TRACING_COLLECTOR;FRONTEND_TRACING_COLLECTOR"`
}

type Logging struct {
	Level  string `yaml:"level,omitempty" env:"OCIS_LOG_LEVEL;FRONTEND_LOG_LEVEL" desc:"The log level."`
	Pretty bool   `yaml:"pretty,omitempty" env:"OCIS_LOG_PRETTY;FRONTEND_LOG_PRETTY" desc:"Activates pretty log output."`
	Color  bool   `yaml:"color,omitempty" env:"OCIS_LOG_COLOR;FRONTEND_LOG_COLOR" desc:"Activates colorized log output."`
	File   string `yaml:"file,omitempty" env:"OCIS_LOG_FILE;FRONTEND_LOG_FILE" desc:"The target log file."`
}

type Service struct {
	Name string `yaml:"-"`
}

type Debug struct {
	Addr   string `yaml:"addr,omitempty" env:"FRONTEND_DEBUG_ADDR"`
	Token  string `yaml:"token,omitempty" env:"FRONTEND_DEBUG_TOKEN"`
	Pprof  bool   `yaml:"pprof,omitempty" env:"FRONTEND_DEBUG_PPROF"`
	Zpages bool   `yaml:"zpages,omitempty" env:"FRONTEND_DEBUG_ZPAGES"`
}

type HTTPConfig struct {
	Addr     string `yaml:"addr,omitempty" env:"FRONTEND_HTTP_ADDR" desc:"The address of the http service."`
	Protocol string `yaml:"protocol,omitempty" env:"FRONTEND_HTTP_PROTOCOL" desc:"The transport protocol of the http service."`
	Prefix   string `yaml:"prefix,omitempty"`
}

// Middleware configures reva middlewares.
type Middleware struct {
	Auth Auth `yaml:"auth,omitempty"`
}

// Auth configures reva http auth middleware.
type Auth struct {
	CredentialsByUserAgent map[string]string `yaml:"credentials_by_user_agent,omitempty"`
}

type Archiver struct {
	MaxNumFiles int64  `yaml:"max_num_files,omitempty"`
	MaxSize     int64  `yaml:"max_size,omitempty"`
	Prefix      string `yaml:"-"`
	Insecure    bool   `yaml:"insecure,omitempty" env:"OCIS_INSECURE;FRONTEND_ARCHIVER_INSECURE"`
}

type AppProvider struct {
	ExternalAddr string `yaml:"external_addr,omitempty"`
	Driver       string `yaml:"driver,omitempty"`
	// WopiDriver   WopiDriver `yaml:"wopi_driver"`
	AppsURL  string `yaml:"-"`
	OpenURL  string `yaml:"-"`
	NewURL   string `yaml:"-"`
	Prefix   string `yaml:"-"`
	Insecure bool   `yaml:"insecure,omitempty" env:"OCIS_INSECURE;FRONTEND_APPPROVIDER_INSECURE"`
}

type DataGateway struct {
	Prefix string
}

type OCS struct {
	Prefix                  string `yaml:"prefix"`
	SharePrefix             string `yaml:"share_prefix"`
	HomeNamespace           string `yaml:"home_namespace"`
	AdditionalInfoAttribute string `yaml:"additional_info_attribute"`
	ResourceInfoCacheTTL    int    `yaml:"resource_info_cache_ttl"`
	CacheWarmupDriver       string `yaml:"cache_warmup_driver"`
	CacheWarmupDrivers      CacheWarmupDrivers
}

type CacheWarmupDrivers struct {
	CBOX CBOXDriver
}

type CBOXDriver struct {
	DBUsername string
	DBPassword string
	DBHost     string
	DBPort     int
	DBName     string
	Namespace  string
}

type AuthMachine struct {
	APIKey string `env:"OCIS_MACHINE_AUTH_API_KEY"`
}

type Checksums struct {
	SupportedTypes      []string `yaml:"supported_types"`
	PreferredUploadType string   `yaml:"preferred_upload_type"`
}
