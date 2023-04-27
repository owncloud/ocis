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
	TransferSecret             string `yaml:"transfer_secret" env:"OCIS_TRANSFER_SECRET;STORAGE_TRANSFER_SECRET" desc:"The storage transfer secret." deprecationVersion:"3.0" removalVersion:"4.0.0" deprecationInfo:"STORAGE_TRANSFER_SECRET changing name for consistency" deprecationReplacement:"OCIS_TRANSFER_SECRET"`
	TransferExpires            int    `yaml:"transfer_expires" env:"GATEWAY_TRANSFER_EXPIRES" desc:"Expiry for the gateway tokens."`
	Cache                      Cache  `yaml:"cache"`

	FrontendPublicURL string `yaml:"frontend_public_url" env:"OCIS_URL;GATEWAY_FRONTEND_PUBLIC_URL" desc:"The public facing URL of the oCIS frontend."`

	UsersEndpoint             string `yaml:"users_endpoint" env:"GATEWAY_USERS_ENDPOINT" desc:"The USERS API endpoint."`
	GroupsEndpoint            string `yaml:"groups_endpoint" env:"GATEWAY_GROUPS_ENDPOINT" desc:"The GROUPS API endpoint."`
	PermissionsEndpoint       string `yaml:"permissions_endpoint" env:"GATEWAY_PERMISSIONS_ENDPOINT" desc:"The SETTINGS API endpoint."`
	SharingEndpoint           string `yaml:"sharing_endpoint" env:"GATEWAY_SHARING_ENDPOINT" desc:"The SHARE API endpoint."`
	AuthBasicEndpoint         string `yaml:"auth_basic_endpoint" env:"GATEWAY_AUTH_BASIC_ENDPOINT" desc:"The AUTH BASIC API endpoint."`
	AuthBearerEndpoint        string `yaml:"auth_bearer_endpoint" env:"GATEWAY_AUTH_BEARER_ENDPOINT" desc:"The AUTH BEARER API endpoint."`
	AuthMachineEndpoint       string `yaml:"auth_machine_endpoint" env:"GATEWAY_AUTH_MACHINE_ENDPOINT" desc:"The AUTH MACHINE API endpoint."`
	StoragePublicLinkEndpoint string `yaml:"storage_public_link_endpoint" env:"GATEWAY_STORAGE_PUBLIC_LINK_ENDPOINT" desc:"The STORAGE PUBLICLINK API endpoint."`
	StorageUsersEndpoint      string `yaml:"storage_users_endpoint" env:"GATEWAY_STORAGE_USERS_ENDPOINT" desc:"The STORAGE USERS API endpoint."`
	StorageSharesEndpoint     string `yaml:"storage_shares_endpoint" env:"GATEWAY_STORAGE_SHARES_ENDPOINT" desc:"The STORAGE SHARES API endpoint."`
	AppRegistryEndpoint       string `yaml:"app_registry_endpoint" env:"GATEWAY_APP_REGISTRY_ENDPOINT" desc:"The APP REGISTRY API endpoint."`

	StorageRegistry StorageRegistry `yaml:"storage_registry"` //TODO: should we even support switching this?

	Supervised bool            `yaml:"-"`
	Context    context.Context `yaml:"-"`
}
type Tracing struct {
	Enabled   bool   `yaml:"enabled" env:"OCIS_TRACING_ENABLED;GATEWAY_TRACING_ENABLED" desc:"Activates tracing."`
	Type      string `yaml:"type" env:"OCIS_TRACING_TYPE;GATEWAY_TRACING_TYPE" desc:"The type of tracing. Defaults to \"\", which is the same as \"jaeger\". Allowed tracing types are \"jaeger\" and \"\" as of now."`
	Endpoint  string `yaml:"endpoint" env:"OCIS_TRACING_ENDPOINT;GATEWAY_TRACING_ENDPOINT" desc:"The endpoint of the tracing agent."`
	Collector string `yaml:"collector" env:"OCIS_TRACING_COLLECTOR;GATEWAY_TRACING_COLLECTOR" desc:"The HTTP endpoint for sending spans directly to a collector, i.e. http://jaeger-collector:14268/api/traces. Only used if the tracing endpoint is unset."`
}

type Log struct {
	Level  string `yaml:"level" env:"OCIS_LOG_LEVEL;GATEWAY_LOG_LEVEL" desc:"The log level. Valid values are: \"panic\", \"fatal\", \"error\", \"warn\", \"info\", \"debug\", \"trace\"."`
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
	Driver              string   `yaml:"driver"` //TODO: configure via env?
	Rules               []string `yaml:"rules"`  //TODO: configure via env?
	JSON                string   `yaml:"json"`   //TODO: configure via env?
	StorageUsersMountID string   `yaml:"storage_users_mount_id" env:"GATEWAY_STORAGE_USERS_MOUNT_ID" desc:"Mount ID of this storage. Admins can set the ID for the storage in this config option manually which is then used to reference the storage. Any reasonable long string is possible, preferably this would be an UUIDv4 format."`
}

// Cache holds cache config
type Cache struct {
	StatCacheStore          string        `yaml:"stat_cache_store" env:"OCIS_CACHE_STORE;GATEWAY_STAT_CACHE_STORE" desc:"The type of the cache store. Supported values are: 'memory', 'ocmem', 'etcd', 'redis', 'redis-sentinel', 'nats-js', 'noop'. See the text description for details."`
	StatCacheNodes          []string      `yaml:"stat_cache_nodes" env:"OCIS_CACHE_STORE_NODES;GATEWAY_STAT_CACHE_STORE_NODES" desc:"A comma separated list of nodes to access the configured store. This has no effect when 'memory' or 'ocmem' stores are configured. Note that the behaviour how nodes are used is dependent on the library of the configured store."`
	StatCacheDatabase       string        `yaml:"stat_cache_database" env:"OCIS_CACHE_DATABASE" desc:"The database name the configured store should use."`
	StatCacheTTL            time.Duration `yaml:"stat_cache_ttl" env:"OCIS_CACHE_TTL;GATEWAY_STAT_CACHE_TTL" desc:"Default time to live for user info in the cache. Only applied when access tokens has no expiration. The duration can be set as number followed by a unit identifier like s, m or h. Defaults to '300s' (300 seconds)."`
	StatCacheSize           int           `yaml:"stat_cache_size" env:"OCIS_CACHE_SIZE;GATEWAY_STAT_CACHE_SIZE" desc:"The maximum quantity of items in the cache. Only applies when store type 'ocmem' is configured. Defaults to 512."`
	ProviderCacheStore      string        `yaml:"provider_cache_store" env:"OCIS_CACHE_STORE;GATEWAY_PROVIDER_CACHE_STORE" desc:"The type of the cache store. Supported values are: 'memory', 'ocmem', 'etcd', 'redis', 'redis-sentinel', 'nats-js', 'noop'. See the text description for details."`
	ProviderCacheNodes      []string      `yaml:"provider_cache_nodes" env:"OCIS_CACHE_STORE_NODES;GATEWAY_PROVIDER_CACHE_STORE_NODES" desc:"A comma separated list of nodes to access the configured store. This has no effect when 'memory' or 'ocmem' stores are configured. Note that the behaviour how nodes are used is dependent on the library of the configured store."`
	ProviderCacheDatabase   string        `yaml:"provider_cache_database" env:"OCIS_CACHE_DATABASE" desc:"The database name the configured store should use."`
	ProviderCacheTTL        time.Duration `yaml:"provider_cache_ttl" env:"OCIS_CACHE_TTL;GATEWAY_PROVIDER_CACHE_TTL" desc:"Default time to live for user info in the cache. Only applied when access tokens has no expiration. The duration can be set as number followed by a unit identifier like s, m or h. Defaults to '300s' (300 seconds)."`
	ProviderCacheSize       int           `yaml:"provider_cache_size" env:"OCIS_CACHE_SIZE;GATEWAY_PROVIDER_CACHE_SIZE" desc:"The maximum quantity of items in the cache. Only applies when store type 'ocmem' is configured. Defaults to 512."`
	CreateHomeCacheStore    string        `yaml:"create_home_cache_store" env:"OCIS_CACHE_STORE;GATEWAY_CREATE_HOME_CACHE_STORE" desc:"The type of the cache store. Supported values are: 'memory', 'ocmem', 'etcd', 'redis', 'redis-sentinel', 'nats-js', 'noop'. See the text description for details."`
	CreateHomeCacheNodes    []string      `yaml:"create_home_cache_nodes" env:"OCIS_CACHE_STORE_NODES;GATEWAY_CREATE_HOME_CACHE_STORE_NODES" desc:"A comma separated list of nodes to access the configured store. This has no effect when 'memory' or 'ocmem' stores are configured. Note that the behaviour how nodes are used is dependent on the library of the configured store."`
	CreateHomeCacheDatabase string        `yaml:"create_home_cache_database" env:"OCIS_CACHE_DATABASE" desc:"The database name the configured store should use."`
	CreateHomeCacheTTL      time.Duration `yaml:"create_home_cache_ttl" env:"OCIS_CACHE_TTL;GATEWAY_CREATE_HOME_CACHE_TTL" desc:"Default time to live for user info in the cache. Only applied when access tokens has no expiration. The duration can be set as number followed by a unit identifier like s, m or h. Defaults to '300s' (300 seconds)."`
	CreateHomeCacheSize     int           `yaml:"create_home_cache_size" env:"OCIS_CACHE_SIZE;GATEWAY_CREATE_HOME_CACHE_SIZE" desc:"The maximum quantity of items in the cache. Only applies when store type 'ocmem' is configured. Defaults to 512."`
}
