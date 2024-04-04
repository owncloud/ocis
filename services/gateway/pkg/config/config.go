package config

import (
	"context"
	"time"

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

	SkipUserGroupsInToken bool `yaml:"skip_user_groups_in_token" env:"GATEWAY_SKIP_USER_GROUPS_IN_TOKEN" desc:"Disables the loading of user's group memberships from the reva access token." introductionVersion:"pre5.0"`

	CommitShareToStorageGrant  bool   `yaml:"commit_share_to_storage_grant" env:"GATEWAY_COMMIT_SHARE_TO_STORAGE_GRANT" desc:"Commit shares to storage grants. This grants access to shared resources for the share receiver directly on the storage." introductionVersion:"pre5.0"`
	ShareFolder                string `yaml:"share_folder_name" env:"GATEWAY_SHARE_FOLDER_NAME" desc:"Name of the share folder in users' home space." introductionVersion:"pre5.0"`
	DisableHomeCreationOnLogin bool   `yaml:"disable_home_creation_on_login" env:"GATEWAY_DISABLE_HOME_CREATION_ON_LOGIN" desc:"Disable creation of the home space on login." introductionVersion:"pre5.0"`
	TransferSecret             string `yaml:"transfer_secret" env:"OCIS_TRANSFER_SECRET" desc:"The storage transfer secret." introductionVersion:"pre5.0"`
	TransferExpires            int    `yaml:"transfer_expires" env:"GATEWAY_TRANSFER_EXPIRES" desc:"Expiry for the gateway tokens." introductionVersion:"pre5.0"`
	Cache                      Cache  `yaml:"cache"`

	FrontendPublicURL string `yaml:"frontend_public_url" env:"OCIS_URL;GATEWAY_FRONTEND_PUBLIC_URL" desc:"The public facing URL of the oCIS frontend." introductionVersion:"pre5.0"`

	UsersEndpoint             string `yaml:"-"`
	GroupsEndpoint            string `yaml:"-"`
	PermissionsEndpoint       string `yaml:"-"`
	SharingEndpoint           string `yaml:"-"`
	AuthBasicEndpoint         string `yaml:"-"`
	AuthBearerEndpoint        string `yaml:"-"`
	AuthMachineEndpoint       string `yaml:"-"`
	AuthServiceEndpoint       string `yaml:"-"`
	StoragePublicLinkEndpoint string `yaml:"-"`
	StorageUsersEndpoint      string `yaml:"-"`
	StorageSharesEndpoint     string `yaml:"-"`
	AppRegistryEndpoint       string `yaml:"-"`
	OCMEndpoint               string `yaml:"-"`

	StorageRegistry StorageRegistry `yaml:"storage_registry"` // TODO: should we even support switching this?

	Supervised bool            `yaml:"-"`
	Context    context.Context `yaml:"-"`
}

type Log struct {
	Level  string `yaml:"level" env:"OCIS_LOG_LEVEL;GATEWAY_LOG_LEVEL" desc:"The log level. Valid values are: 'panic', 'fatal', 'error', 'warn', 'info', 'debug', 'trace'." introductionVersion:"pre5.0"`
	Pretty bool   `yaml:"pretty" env:"OCIS_LOG_PRETTY;GATEWAY_LOG_PRETTY" desc:"Activates pretty log output." introductionVersion:"pre5.0"`
	Color  bool   `yaml:"color" env:"OCIS_LOG_COLOR;GATEWAY_LOG_COLOR" desc:"Activates colorized log output." introductionVersion:"pre5.0"`
	File   string `yaml:"file" env:"OCIS_LOG_FILE;GATEWAY_LOG_FILE" desc:"The path to the log file. Activates logging to this file if set." introductionVersion:"pre5.0"`
}

type Service struct {
	Name string `yaml:"-"`
}

type Debug struct {
	Addr   string `yaml:"addr" env:"GATEWAY_DEBUG_ADDR" desc:"Bind address of the debug server, where metrics, health, config and debug endpoints will be exposed." introductionVersion:"pre5.0"`
	Token  string `yaml:"token" env:"GATEWAY_DEBUG_TOKEN" desc:"Token to secure the metrics endpoint." introductionVersion:"pre5.0"`
	Pprof  bool   `yaml:"pprof" env:"GATEWAY_DEBUG_PPROF" desc:"Enables pprof, which can be used for profiling." introductionVersion:"pre5.0"`
	Zpages bool   `yaml:"zpages" env:"GATEWAY_DEBUG_ZPAGES" desc:"Enables zpages, which can be used for collecting and viewing in-memory traces." introductionVersion:"pre5.0"`
}

type GRPCConfig struct {
	Addr      string                 `yaml:"addr" env:"OCIS_GATEWAY_GRPC_ADDR;GATEWAY_GRPC_ADDR" desc:"The bind address of the GRPC service." introductionVersion:"pre5.0"`
	TLS       *shared.GRPCServiceTLS `yaml:"tls"`
	Namespace string                 `yaml:"-"`
	Protocol  string                 `yaml:"protocol" env:"GATEWAY_GRPC_PROTOCOL" desc:"The transport protocol of the GRPC service." introductionVersion:"pre5.0"`
}

type StorageRegistry struct {
	Driver              string   `yaml:"driver" env:"GATEWAY_STORAGE_REGISTRY_DRIVER" desc:"The driver name of the storage registry to use." introductionVersion:"pre5.0"`
	Rules               []string `yaml:"rules" env:"GATEWAY_STORAGE_REGISTRY_RULES" desc:"The rules for the storage registry. See the Environment Variable Types description for more details." introductionVersion:"pre5.0"`
	JSON                string   `yaml:"json" env:"GATEWAY_STORAGE_REGISTRY_CONFIG_JSON" desc:"Additional configuration for the storage registry in json format." introductionVersion:"pre5.0"`
	StorageUsersMountID string   `yaml:"storage_users_mount_id" env:"GATEWAY_STORAGE_USERS_MOUNT_ID" desc:"Mount ID of this storage. Admins can set the ID for the storage in this config option manually which is then used to reference the storage. Any reasonable long string is possible, preferably this would be an UUIDv4 format." introductionVersion:"pre5.0"`
}

// Cache holds cache config
type Cache struct {
	ProviderCacheStore                string        `yaml:"provider_cache_store" env:"OCIS_CACHE_STORE;GATEWAY_PROVIDER_CACHE_STORE" desc:"The type of the cache store. Supported values are: 'memory', 'redis-sentinel', 'nats-js-kv', 'noop'. See the text description for details." introductionVersion:"pre5.0"`
	ProviderCacheNodes                []string      `yaml:"provider_cache_nodes" env:"OCIS_CACHE_STORE_NODES;GATEWAY_PROVIDER_CACHE_STORE_NODES" desc:"A list of nodes to access the configured store. This has no effect when 'memory' or 'ocmem' stores are configured. Note that the behaviour how nodes are used is dependent on the library of the configured store. See the Environment Variable Types description for more details." introductionVersion:"pre5.0"`
	ProviderCacheDatabase             string        `yaml:"provider_cache_database" env:"OCIS_CACHE_DATABASE" desc:"The database name the configured store should use." introductionVersion:"pre5.0"`
	ProviderCacheTTL                  time.Duration `yaml:"provider_cache_ttl" env:"OCIS_CACHE_TTL;GATEWAY_PROVIDER_CACHE_TTL" desc:"Default time to live for user info in the cache. Only applied when access tokens has no expiration. See the Environment Variable Types description for more details." introductionVersion:"pre5.0"`
	ProviderCacheSize                 int           `yaml:"provider_cache_size" env:"OCIS_CACHE_SIZE;GATEWAY_PROVIDER_CACHE_SIZE" desc:"The maximum quantity of items in the cache. Only applies when store type 'ocmem' is configured. Defaults to 512 which is derived from the ocmem package though not explicitly set as default." introductionVersion:"pre5.0"`
	ProviderCacheDisablePersistence   bool          `yaml:"provider_cache_disable_persistence" env:"OCIS_CACHE_DISABLE_PERSISTENCE;GATEWAY_PROVIDER_CACHE_DISABLE_PERSISTENCE" desc:"Disables persistence of the provider cache. Only applies when store type 'nats-js-kv' is configured. Defaults to false." introductionVersion:"5.0"`
	ProviderCacheAuthUsername         string        `yaml:"provider_cache_auth_username" env:"OCIS_CACHE_AUTH_USERNAME;GATEWAY_PROVIDER_CACHE_AUTH_USERNAME" desc:"The username to use for authentication. Only applies when store type 'nats-js-kv' is configured." introductionVersion:"5.0"`
	ProviderCacheAuthPassword         string        `yaml:"provider_cache_auth_password" env:"OCIS_CACHE_AUTH_PASSWORD;GATEWAY_PROVIDER_CACHE_AUTH_PASSWORD" desc:"The password to use for authentication. Only applies when store type 'nats-js-kv' is configured." introductionVersion:"5.0"`
	CreateHomeCacheStore              string        `yaml:"create_home_cache_store" env:"OCIS_CACHE_STORE;GATEWAY_CREATE_HOME_CACHE_STORE" desc:"The type of the cache store. Supported values are: 'memory', 'redis-sentinel', 'nats-js-kv', 'noop'. See the text description for details." introductionVersion:"pre5.0"`
	CreateHomeCacheNodes              []string      `yaml:"create_home_cache_nodes" env:"OCIS_CACHE_STORE_NODES;GATEWAY_CREATE_HOME_CACHE_STORE_NODES" desc:"A list of nodes to access the configured store. This has no effect when 'memory' or 'ocmem' stores are configured. Note that the behaviour how nodes are used is dependent on the library of the configured store. See the Environment Variable Types description for more details." introductionVersion:"pre5.0"`
	CreateHomeCacheDatabase           string        `yaml:"create_home_cache_database" env:"OCIS_CACHE_DATABASE" desc:"The database name the configured store should use." introductionVersion:"pre5.0"`
	CreateHomeCacheTTL                time.Duration `yaml:"create_home_cache_ttl" env:"OCIS_CACHE_TTL;GATEWAY_CREATE_HOME_CACHE_TTL" desc:"Default time to live for user info in the cache. Only applied when access tokens has no expiration. See the Environment Variable Types description for more details." introductionVersion:"pre5.0"`
	CreateHomeCacheSize               int           `yaml:"create_home_cache_size" env:"OCIS_CACHE_SIZE;GATEWAY_CREATE_HOME_CACHE_SIZE" desc:"The maximum quantity of items in the cache. Only applies when store type 'ocmem' is configured. Defaults to 512 which is derived from the ocmem package though not explicitly set as default." introductionVersion:"pre5.0"`
	CreateHomeCacheDisablePersistence bool          `yaml:"create_home_cache_disable_persistence" env:"OCIS_CACHE_DISABLE_PERSISTENCE;GATEWAY_CREATE_HOME_CACHE_DISABLE_PERSISTENCE" desc:"Disables persistence of the create home cache. Only applies when store type 'nats-js-kv' is configured. Defaults to false." introductionVersion:"5.0"`
	CreateHomeCacheAuthUsername       string        `yaml:"create_home_cache_auth_username" env:"OCIS_CACHE_AUTH_USERNAME;GATEWAY_CREATE_HOME_CACHE_AUTH_USERNAME" desc:"The username to use for authentication. Only applies when store type 'nats-js-kv' is configured." introductionVersion:"5.0"`
	CreateHomeCacheAuthPassword       string        `yaml:"create_home_cache_auth_password" env:"OCIS_CACHE_AUTH_PASSWORD;GATEWAY_CREATE_HOME_CACHE_AUTH_PASSWORD" desc:"The password to use for authentication. Only applies when store type 'nats-js-kv' is configured." introductionVersion:"5.0"`
}
