package config

import "github.com/owncloud/ocis/ocis-pkg/shared"

type Config struct {
	*shared.Commons `yaml:"-"`
	Service         Service  `yaml:"-"`
	Tracing         *Tracing `yaml:"tracing"`
	Logging         *Logging `yaml:"log"`
	Debug           Debug    `yaml:"debug"`
	Supervised      bool

	HTTP HTTPConfig `yaml:"http"`

	// JWTSecret used to verify reva access token

	TransferSecret string `yaml:"transfer_secret" env:"STORAGE_TRANSFER_SECRET"`

	JWTSecret string `yaml:"jwt_secret"`
	GatewayEndpoint       string
	SkipUserGroupsInToken bool

	EnableFavorites          bool `yaml:"favorites"`
	EnableProjectSpaces      bool
	UploadMaxChunkSize       int    `yaml:"upload_max_chunk_size"`
	UploadHTTPMethodOverride string `yaml:"upload_http_method_override"`
	DefaultUploadProtocol    string `yaml:"default_upload_protocol"`

	PublicURL string `yaml:"public_url" env:"OCIS_URL;FRONTEND_PUBLIC_URL"`

	Archiver    Archiver
	AppProvider AppProvider
	DataGateway DataGateway
	OCS         OCS
	AuthMachine AuthMachine
	Checksums   Checksums

	Middleware Middleware
}
type Tracing struct {
	Enabled   bool   `yaml:"enabled" env:"OCIS_TRACING_ENABLED;FRONTEND_TRACING_ENABLED" desc:"Activates tracing."`
	Type      string `yaml:"type" env:"OCIS_TRACING_TYPE;FRONTEND_TRACING_TYPE"`
	Endpoint  string `yaml:"endpoint" env:"OCIS_TRACING_ENDPOINT;FRONTEND_TRACING_ENDPOINT" desc:"The endpoint to the tracing collector."`
	Collector string `yaml:"collector" env:"OCIS_TRACING_COLLECTOR;FRONTEND_TRACING_COLLECTOR"`
}

type Logging struct {
	Level  string `yaml:"level" env:"OCIS_LOG_LEVEL;FRONTEND_LOG_LEVEL" desc:"The log level."`
	Pretty bool   `yaml:"pretty" env:"OCIS_LOG_PRETTY;FRONTEND_LOG_PRETTY" desc:"Activates pretty log output."`
	Color  bool   `yaml:"color" env:"OCIS_LOG_COLOR;FRONTEND_LOG_COLOR" desc:"Activates colorized log output."`
	File   string `yaml:"file" env:"OCIS_LOG_FILE;FRONTEND_LOG_FILE" desc:"The target log file."`
}

type Service struct {
	Name string `yaml:"-"`
}

type Debug struct {
	Addr   string `yaml:"addr" env:"FRONTEND_DEBUG_ADDR"`
	Token  string `yaml:"token" env:"FRONTEND_DEBUG_TOKEN"`
	Pprof  bool   `yaml:"pprof" env:"FRONTEND_DEBUG_PPROF"`
	Zpages bool   `yaml:"zpages" env:"FRONTEND_DEBUG_ZPAGES"`
}

type HTTPConfig struct {
	Addr     string `yaml:"addr" env:"FRONTEND_HTTP_ADDR" desc:"The address of the http service."`
	Protocol string `yaml:"protocol" env:"FRONTEND_HTTP_PROTOCOL" desc:"The transport protocol of the http service."`
	Prefix   string `yaml:"prefix"`
}

// Middleware configures reva middlewares.
type Middleware struct {
	Auth Auth `yaml:"auth"`
}

// Auth configures reva http auth middleware.
type Auth struct {
	CredentialsByUserAgent map[string]string `yaml:"credentials_by_user_agenr"`
}

type Archiver struct {
	MaxNumFiles int64 `yaml:"max_num_files"`
	MaxSize     int64 `yaml:"max_size"`
	Prefix      string
	Insecure    bool `env:"OCIS_INSECURE;FRONTEND_ARCHIVER_INSECURE"`
}

type AppProvider struct {
	ExternalAddr string `yaml:"external_addr"`
	Driver       string `yaml:"driver"`
	// WopiDriver   WopiDriver `yaml:"wopi_driver"`
	AppsURL  string `yaml:"apps_url"`
	OpenURL  string `yaml:"open_url"`
	NewURL   string `yaml:"new_url"`
	Prefix   string
	Insecure bool `env:"OCIS_INSECURE;FRONTEND_APPPROVIDER_INSECURE"`
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
