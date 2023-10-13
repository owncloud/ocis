package config

import (
	"context"
	"time"

	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
	settingsmsg "github.com/owncloud/ocis/v2/protogen/gen/ocis/messages/settings/v0"
	"go-micro.dev/v4/client"
)

// Config combines all available configuration parts.
type Config struct {
	Commons *shared.Commons `yaml:"-"` // don't use this directly as configuration for a service

	Service Service `yaml:"-"`

	Tracing *Tracing `yaml:"tracing"`
	Log     *Log     `yaml:"log"`
	Debug   Debug    `yaml:"debug"`

	HTTP HTTP       `yaml:"http"`
	GRPC GRPCConfig `yaml:"grpc"`

	GRPCClientTLS *shared.GRPCClientTLS `yaml:"grpc_client_tls"`
	GrpcClient    client.Client         `yaml:"-"`

	StoreType   string                `yaml:"store_type" env:"SETTINGS_STORE_TYPE" desc:"Store type configures the persistency driver. Supported values are 'metadata' and 'filesystem'. Note that the value 'filesystem' is considered deprecated."`
	DataPath    string                `yaml:"data_path" env:"SETTINGS_DATA_PATH" desc:"The directory where the filesystem storage will store ocis settings. If not defined, the root directory derives from $OCIS_BASE_DATA_PATH:/settings."`
	Metadata    Metadata              `yaml:"metadata_config"`
	BundlesPath string                `yaml:"bundles_path" env:"SETTINGS_BUNDLES_PATH" desc:"The path to a JSON file with a list of bundles. If not defined, the default bundles will be loaded."`
	Bundles     []*settingsmsg.Bundle `yaml:"-"`

	AdminUserID string `yaml:"admin_user_id" env:"OCIS_ADMIN_USER_ID;SETTINGS_ADMIN_USER_ID" desc:"ID of the user that should receive admin privileges. Consider that the UUID can be encoded in some LDAP deployment configurations like in .ldif files. These need to be decoded beforehand."`

	TokenManager *TokenManager `yaml:"token_manager"`

	SetupDefaultAssignments bool `yaml:"set_default_assignments" env:"SETTINGS_SETUP_DEFAULT_ASSIGNMENTS;IDM_CREATE_DEMO_USERS" desc:"The default role assignments the demo users should be setup."`

	ServiceAccountIDAdmin string `yaml:"service_account_id_admin" env:"OCIS_SERVICE_ACCOUNT_ID;SETTINGS_SERVICE_ACCOUNT_ID_ADMIN" desc:"The ID of the service account having the admin role. See the 'auth-service' service description for more details."`

	DefaultLanguage string `yaml:"default_language" env:"OCIS_DEFAULT_LANGUAGE" desc:"(optional) The default language. If not defined, English will be used as default. See the documentation for more details."`

	Context context.Context `yaml:"-"`
}

// Metadata configures the metadata store to use
type Metadata struct {
	GatewayAddress string `yaml:"gateway_addr" env:"STORAGE_GATEWAY_GRPC_ADDR" desc:"GRPC address of the STORAGE-SYSTEM service."`
	StorageAddress string `yaml:"storage_addr" env:"STORAGE_GRPC_ADDR" desc:"GRPC address of the STORAGE-SYSTEM service."`

	SystemUserID     string `yaml:"system_user_id" env:"OCIS_SYSTEM_USER_ID;SETTINGS_SYSTEM_USER_ID" desc:"ID of the oCIS STORAGE-SYSTEM system user. Admins need to set the ID for the STORAGE-SYSTEM system user in this config option which is then used to reference the user. Any reasonable long string is possible, preferably this would be an UUIDv4 format."`
	SystemUserIDP    string `yaml:"system_user_idp" env:"OCIS_SYSTEM_USER_IDP;SETTINGS_SYSTEM_USER_IDP" desc:"IDP of the oCIS STORAGE-SYSTEM system user."`
	SystemUserAPIKey string `yaml:"system_user_api_key" env:"OCIS_SYSTEM_USER_API_KEY" desc:"API key for the STORAGE-SYSTEM system user."`
	Cache            *Cache `yaml:"cache"`
}

// Cache configures the cache of the Metadata store
type Cache struct {
	Store          string        `yaml:"store" env:"OCIS_CACHE_STORE;SETTINGS_CACHE_STORE" desc:"The type of the cache store. Supported values are: 'memory', 'ocmem', 'etcd', 'redis', 'redis-sentinel', 'nats-js', 'noop'. See the text description for details."`
	Nodes          []string      `yaml:"addresses" env:"OCIS_CACHE_STORE_NODES;SETTINGS_CACHE_STORE_NODES" desc:"A comma separated list of nodes to access the configured store. This has no effect when 'memory' or 'ocmem' stores are configured. Note that the behaviour how nodes are used is dependent on the library of the configured store."`
	Database       string        `yaml:"database" env:"OCIS_CACHE_DATABASE" desc:"The database name the configured store should use."`
	FileTable      string        `yaml:"files_table" env:"SETTINGS_FILE_CACHE_TABLE" desc:"The database table the store should use for the file cache."`
	DirectoryTable string        `yaml:"directories_table" env:"SETTINGS_DIRECTORY_CACHE_TABLE" desc:"The database table the store should use for the directory cache."`
	TTL            time.Duration `yaml:"ttl" env:"OCIS_CACHE_TTL;SETTINGS_CACHE_TTL" desc:"Default time to live for entries in the cache. Only applied when access tokens has no expiration. The duration can be set as number followed by a unit identifier like s, m or h. Defaults to '10m' (10 minutes)."`
	Size           int           `yaml:"size" env:"OCIS_CACHE_SIZE;SETTINGS_CACHE_SIZE" desc:"The maximum quantity of items in the cache. Only applies when store type 'ocmem' is configured. Defaults to 512."`
}
