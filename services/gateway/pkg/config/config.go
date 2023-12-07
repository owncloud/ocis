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

	SkipUserGroupsInToken bool `yaml:"skip_user_groups_in_token" env:"GATEWAY_SKIP_USER_GROUPS_IN_TOKEN" desc:"Disables the loading of user's group memberships from the reva access token."`

	CommitShareToStorageGrant  bool   `yaml:"commit_share_to_storage_grant" env:"GATEWAY_COMMIT_SHARE_TO_STORAGE_GRANT" desc:"Commit shares to storage grants. This grants access to shared resources for the share receiver directly on the storage."`
	ShareFolder                string `yaml:"share_folder_name" env:"GATEWAY_SHARE_FOLDER_NAME" desc:"Name of the share folder in users' home space."`
	DisableHomeCreationOnLogin bool   `yaml:"disable_home_creation_on_login" env:"GATEWAY_DISABLE_HOME_CREATION_ON_LOGIN" desc:"Disable creation of the home space on login."`
	TransferSecret             string `yaml:"transfer_secret" env:"OCIS_TRANSFER_SECRET" desc:"The storage transfer secret."`
	TransferExpires            int    `yaml:"transfer_expires" env:"GATEWAY_TRANSFER_EXPIRES" desc:"Expiry for the gateway tokens."`
	Cache                      Cache  `yaml:"cache"`

	FrontendPublicURL string `yaml:"frontend_public_url" env:"OCIS_URL;GATEWAY_FRONTEND_PUBLIC_URL" desc:"The public facing URL of the oCIS frontend."`

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
	Level  string `yaml:"level" env:"OCIS_LOG_LEVEL;GATEWAY_LOG_LEVEL" desc:"The log level. Valid values are: 'panic', 'fatal', 'error', 'warn', 'info', 'debug', 'trace'."`
	Pretty bool   `yaml:"pretty" env:"OCIS_LOG_PRETTY;GATEWAY_LOG_PRETTY" desc:"Activates pretty log output."`
	Color  bool   `yaml:"color" env:"OCIS_LOG_COLOR;GATEWAY_LOG_COLOR" desc:"Activates colorized log output."`
	File   string `yaml:"file" env:"OCIS_LOG_FILE;GATEWAY_LOG_FILE" desc:"The path to the log file. Activates logging to this file if set."`
}

type Service struct {
	Name string `yaml:"-"`
}

type Debug struct {
	Addr   string `yaml:"addr" env:"GATEWAY_DEBUG_ADDR" desc:"Bind address of the debug server, where metrics, health, config and debug endpoints will be exposed."`
	Token  string `yaml:"token" env:"GATEWAY_DEBUG_TOKEN" desc:"Token to secure the metrics endpoint."`
	Pprof  bool   `yaml:"pprof" env:"GATEWAY_DEBUG_PPROF" desc:"Enables pprof, which can be used for profiling."`
	Zpages bool   `yaml:"zpages" env:"GATEWAY_DEBUG_ZPAGES" desc:"Enables zpages, which can be used for collecting and viewing in-memory traces."`
}

type GRPCConfig struct {
	Addr      string                 `yaml:"addr" env:"GATEWAY_GRPC_ADDR" desc:"The bind address of the GRPC service."`
	TLS       *shared.GRPCServiceTLS `yaml:"tls"`
	Namespace string                 `yaml:"-"`
	Protocol  string                 `yaml:"protocol" env:"GATEWAY_GRPC_PROTOCOL" desc:"The transport protocol of the GRPC service."`
}

type StorageRegistry struct {
	Driver              string   `yaml:"driver" env:"GATEWAY_STORAGE_REGISTRY_DRIVER" desc:"The driver name of the storage registry to use."`
	Rules               []string `yaml:"rules" env:"GATEWAY_STORAGE_REGISTRY_RULES" desc:"The rules for the storage registry. See the Environment Variable Types description for more details."`
	JSON                string   `yaml:"json" env:"GATEWAY_STORAGE_REGISTRY_CONFIG_JSON" desc:"Additional configuration for the storage registry in json format."`
	StorageUsersMountID string   `yaml:"storage_users_mount_id" env:"GATEWAY_STORAGE_USERS_MOUNT_ID" desc:"Mount ID of this storage. Admins can set the ID for the storage in this config option manually which is then used to reference the storage. Any reasonable long string is possible, preferably this would be an UUIDv4 format."`
}

// Cache holds cache config
type Cache struct {
	StatCacheStore          string        // NOTE: The stat cache is not working atm. Hence we block configuring it
	StatCacheNodes          []string      `yaml:"stat_cache_nodes" env:"OCIS_CACHE_STORE_NODES;GATEWAY_STAT_CACHE_STORE_NODES" desc:"A list of nodes to access the configured store. This has no effect when 'memory' or 'ocmem' stores are configured. Note that the behaviour how nodes are used is dependent on the library of the configured store. See the Environment Variable Types description for more details."`
	StatCacheDatabase       string        `yaml:"stat_cache_database" env:"OCIS_CACHE_DATABASE" desc:"The database name the configured store should use."`
	StatCacheTTL            time.Duration `yaml:"stat_cache_ttl" env:"OCIS_CACHE_TTL;GATEWAY_STAT_CACHE_TTL" desc:"Default time to live for user info in the cache. Only applied when access tokens has no expiration. See the Environment Variable Types description for more details."`
	StatCacheSize           int           `yaml:"stat_cache_size" env:"OCIS_CACHE_SIZE;GATEWAY_STAT_CACHE_SIZE" desc:"The maximum quantity of items in the cache. Only applies when store type 'ocmem' is configured. Defaults to 512 which is derived from the ocmem package though not exclicitely set as default."`
	ProviderCacheStore      string        `yaml:"provider_cache_store" env:"OCIS_CACHE_STORE;GATEWAY_PROVIDER_CACHE_STORE" desc:"The type of the cache store. Supported values are: 'memory', 'ocmem', 'etcd', 'redis', 'redis-sentinel', 'nats-js', 'noop'. See the text description for details."`
	ProviderCacheNodes      []string      `yaml:"provider_cache_nodes" env:"OCIS_CACHE_STORE_NODES;GATEWAY_PROVIDER_CACHE_STORE_NODES" desc:"A list of nodes to access the configured store. This has no effect when 'memory' or 'ocmem' stores are configured. Note that the behaviour how nodes are used is dependent on the library of the configured store. See the Environment Variable Types description for more details."`
	ProviderCacheDatabase   string        `yaml:"provider_cache_database" env:"OCIS_CACHE_DATABASE" desc:"The database name the configured store should use."`
	ProviderCacheTTL        time.Duration `yaml:"provider_cache_ttl" env:"OCIS_CACHE_TTL;GATEWAY_PROVIDER_CACHE_TTL" desc:"Default time to live for user info in the cache. Only applied when access tokens has no expiration. See the Environment Variable Types description for more details."`
	ProviderCacheSize       int           `yaml:"provider_cache_size" env:"OCIS_CACHE_SIZE;GATEWAY_PROVIDER_CACHE_SIZE" desc:"The maximum quantity of items in the cache. Only applies when store type 'ocmem' is configured. Defaults to 512 which is derived from the ocmem package though not exclicitely set as default."`
	CreateHomeCacheStore    string        `yaml:"create_home_cache_store" env:"OCIS_CACHE_STORE;GATEWAY_CREATE_HOME_CACHE_STORE" desc:"The type of the cache store. Supported values are: 'memory', 'ocmem', 'etcd', 'redis', 'redis-sentinel', 'nats-js', 'noop'. See the text description for details."`
	CreateHomeCacheNodes    []string      `yaml:"create_home_cache_nodes" env:"OCIS_CACHE_STORE_NODES;GATEWAY_CREATE_HOME_CACHE_STORE_NODES" desc:"A list of nodes to access the configured store. This has no effect when 'memory' or 'ocmem' stores are configured. Note that the behaviour how nodes are used is dependent on the library of the configured store. See the Environment Variable Types description for more details."`
	CreateHomeCacheDatabase string        `yaml:"create_home_cache_database" env:"OCIS_CACHE_DATABASE" desc:"The database name the configured store should use."`
	CreateHomeCacheTTL      time.Duration `yaml:"create_home_cache_ttl" env:"OCIS_CACHE_TTL;GATEWAY_CREATE_HOME_CACHE_TTL" desc:"Default time to live for user info in the cache. Only applied when access tokens has no expiration. See the Environment Variable Types description for more details."`
	CreateHomeCacheSize     int           `yaml:"create_home_cache_size" env:"OCIS_CACHE_SIZE;GATEWAY_CREATE_HOME_CACHE_SIZE" desc:"The maximum quantity of items in the cache. Only applies when store type 'ocmem' is configured. Defaults to 512 which is derived from the ocmem package though not exclicitely set as default."`
}
