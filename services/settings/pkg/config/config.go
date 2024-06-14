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

	Metadata    Metadata              `yaml:"metadata_config"`
	BundlesPath string                `yaml:"bundles_path" env:"SETTINGS_BUNDLES_PATH" desc:"The path to a JSON file with a list of bundles. If not defined, the default bundles will be loaded." introductionVersion:"pre5.0"`
	Bundles     []*settingsmsg.Bundle `yaml:"-"`

	AdminUserID string `yaml:"admin_user_id" env:"OCIS_ADMIN_USER_ID;SETTINGS_ADMIN_USER_ID" desc:"ID of the user that should receive admin privileges. Consider that the UUID can be encoded in some LDAP deployment configurations like in .ldif files. These need to be decoded beforehand." introductionVersion:"pre5.0"`

	TokenManager *TokenManager `yaml:"token_manager"`

	SetupDefaultAssignments bool `yaml:"set_default_assignments" env:"SETTINGS_SETUP_DEFAULT_ASSIGNMENTS;IDM_CREATE_DEMO_USERS" desc:"The default role assignments the demo users should be setup." introductionVersion:"pre5.0"`

	ServiceAccountIDs []string `yaml:"service_account_ids" env:"SETTINGS_SERVICE_ACCOUNT_IDS;OCIS_SERVICE_ACCOUNT_ID" desc:"The list of all service account IDs. These will be assigned the hidden 'service-account' role. Note: When using 'OCIS_SERVICE_ACCOUNT_ID' this will contain only one value while 'SETTINGS_SERVICE_ACCOUNT_IDS' can have multiple. See the 'auth-service' service description for more details about service accounts." introductionVersion:"5.0"`

	DefaultLanguage string `yaml:"default_language" env:"OCIS_DEFAULT_LANGUAGE" desc:"The default language used by services and the WebUI. If not defined, English will be used as default. See the documentation for more details." introductionVersion:"5.0"`

	Context context.Context `yaml:"-"`
}

// Metadata configures the metadata store to use
type Metadata struct {
	GatewayAddress string `yaml:"gateway_addr" env:"STORAGE_GATEWAY_GRPC_ADDR" desc:"GRPC address of the STORAGE-SYSTEM service." introductionVersion:"pre5.0"`
	StorageAddress string `yaml:"storage_addr" env:"STORAGE_GRPC_ADDR" desc:"GRPC address of the STORAGE-SYSTEM service." introductionVersion:"pre5.0"`

	SystemUserID     string `yaml:"system_user_id" env:"OCIS_SYSTEM_USER_ID;SETTINGS_SYSTEM_USER_ID" desc:"ID of the oCIS STORAGE-SYSTEM system user. Admins need to set the ID for the STORAGE-SYSTEM system user in this config option which is then used to reference the user. Any reasonable long string is possible, preferably this would be an UUIDv4 format." introductionVersion:"pre5.0"`
	SystemUserIDP    string `yaml:"system_user_idp" env:"OCIS_SYSTEM_USER_IDP;SETTINGS_SYSTEM_USER_IDP" desc:"IDP of the oCIS STORAGE-SYSTEM system user." introductionVersion:"pre5.0"`
	SystemUserAPIKey string `yaml:"system_user_api_key" env:"OCIS_SYSTEM_USER_API_KEY" desc:"API key for the STORAGE-SYSTEM system user." introductionVersion:"pre5.0"`
	Cache            *Cache `yaml:"cache"`
}

// Cache configures the cache of the Metadata store
type Cache struct {
	Store              string        `yaml:"store" env:"OCIS_CACHE_STORE;SETTINGS_CACHE_STORE" desc:"The type of the cache store. Supported values are: 'memory', 'redis-sentinel', 'nats-js-kv', 'noop'. See the text description for details." introductionVersion:"pre5.0"`
	Nodes              []string      `yaml:"addresses" env:"OCIS_CACHE_STORE_NODES;SETTINGS_CACHE_STORE_NODES" desc:"A list of nodes to access the configured store. This has no effect when 'memory' or 'ocmem' stores are configured. Note that the behaviour how nodes are used is dependent on the library of the configured store. See the Environment Variable Types description for more details." introductionVersion:"pre5.0"`
	Database           string        `yaml:"database" env:"OCIS_CACHE_DATABASE" desc:"The database name the configured store should use." introductionVersion:"pre5.0"`
	FileTable          string        `yaml:"files_table" env:"SETTINGS_FILE_CACHE_TABLE" desc:"The database table the store should use for the file cache." introductionVersion:"pre5.0"`
	DirectoryTable     string        `yaml:"directories_table" env:"SETTINGS_DIRECTORY_CACHE_TABLE" desc:"The database table the store should use for the directory cache." introductionVersion:"pre5.0"`
	TTL                time.Duration `yaml:"ttl" env:"OCIS_CACHE_TTL;SETTINGS_CACHE_TTL" desc:"Default time to live for entries in the cache. Only applied when access tokens has no expiration. See the Environment Variable Types description for more details." introductionVersion:"pre5.0"`
	Size               int           `yaml:"size" env:"OCIS_CACHE_SIZE;SETTINGS_CACHE_SIZE" desc:"The maximum quantity of items in the cache. Only applies when store type 'ocmem' is configured. Defaults to 512 which is derived from the ocmem package though not exclicitly set as default." introductionVersion:"pre5.0"`
	DisablePersistence bool          `yaml:"disable_persistence" env:"OCIS_CACHE_DISABLE_PERSISTENCE;SETTINGS_CACHE_DISABLE_PERSISTENCE" desc:"Disables persistence of the cache. Only applies when store type 'nats-js-kv' is configured. Defaults to false." introductionVersion:"5.0"`
	AuthUsername       string        `yaml:"username" env:"OCIS_CACHE_AUTH_USERNAME;SETTINGS_CACHE_AUTH_USERNAME" desc:"The username to authenticate with the cache. Only applies when store type 'nats-js-kv' is configured." introductionVersion:"5.0"`
	AuthPassword       string        `yaml:"password" env:"OCIS_CACHE_AUTH_PASSWORD;SETTINGS_CACHE_AUTH_PASSWORD" desc:"The password to authenticate with the cache. Only applies when store type 'nats-js-kv' is configured." introductionVersion:"5.0"`
}
