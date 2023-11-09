package config

import (
	"context"
	"time"

	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
)

// Config holds Config config
type Config struct {
	Commons *shared.Commons `yaml:"-"` // don't use this directly as configuration for a service
	Service Service         `yaml:"-"`
	Tracing *Tracing        `yaml:"tracing"`
	Log     *Log            `yaml:"log"`
	Debug   Debug           `yaml:"debug"`

	GRPC GRPCConfig `yaml:"grpc"`
	HTTP HTTPConfig `yaml:"http"`

	TokenManager     *TokenManager `yaml:"token_manager"`
	Reva             *shared.Reva  `yaml:"reva"`
	SystemUserID     string        `yaml:"system_user_id" env:"OCIS_SYSTEM_USER_ID" desc:"ID of the oCIS storage-system system user. Admins need to set the ID for the STORAGE-SYSTEM system user in this config option which is then used to reference the user. Any reasonable long string is possible, preferably this would be an UUIDv4 format."`
	SystemUserAPIKey string        `yaml:"system_user_api_key" env:"OCIS_SYSTEM_USER_API_KEY" desc:"API key for the STORAGE-SYSTEM system user."`

	SkipUserGroupsInToken bool `yaml:"skip_user_groups_in_token" env:"STORAGE_SYSTEM_SKIP_USER_GROUPS_IN_TOKEN" desc:"Disables the loading of user's group memberships from the reva access token."`

	FileMetadataCache Cache   `yaml:"cache"`
	Driver            string  `yaml:"driver" env:"STORAGE_SYSTEM_DRIVER" desc:"The driver which should be used by the service."`
	Drivers           Drivers `yaml:"drivers"`
	DataServerURL     string  `yaml:"data_server_url" env:"STORAGE_SYSTEM_DATA_SERVER_URL" desc:"URL of the data server, needs to be reachable by other services using this service."`

	Supervised bool            `yaml:"-"`
	Context    context.Context `yaml:"-"`
}

// Log holds Log config
type Log struct {
	Level  string `yaml:"level" env:"OCIS_LOG_LEVEL;STORAGE_SYSTEM_LOG_LEVEL" desc:"The log level. Valid values are: 'panic', 'fatal', 'error', 'warn', 'info', 'debug', 'trace'."`
	Pretty bool   `yaml:"pretty" env:"OCIS_LOG_PRETTY;STORAGE_SYSTEM_LOG_PRETTY" desc:"Activates pretty log output."`
	Color  bool   `yaml:"color" env:"OCIS_LOG_COLOR;STORAGE_SYSTEM_LOG_COLOR" desc:"Activates colorized log output."`
	File   string `yaml:"file" env:"OCIS_LOG_FILE;STORAGE_SYSTEM_LOG_FILE" desc:"The path to the log file. Activates logging to this file if set."`
}

// Service holds Service config
type Service struct {
	Name string `yaml:"-"`
}

// Debug holds Debug config
type Debug struct {
	Addr   string `yaml:"addr" env:"STORAGE_SYSTEM_DEBUG_ADDR" desc:"Bind address of the debug server, where metrics, health, config and debug endpoints will be exposed."`
	Token  string `yaml:"token" env:"STORAGE_SYSTEM_DEBUG_TOKEN" desc:"Token to secure the metrics endpoint"`
	Pprof  bool   `yaml:"pprof" env:"STORAGE_SYSTEM_DEBUG_PPROF" desc:"Enables pprof, which can be used for profiling"`
	Zpages bool   `yaml:"zpages" env:"STORAGE_SYSTEM_DEBUG_ZPAGES" desc:"Enables zpages, which can be used for collecting and viewing in-memory traces."`
}

// GRPCConfig holds GRPCConfig config
type GRPCConfig struct {
	Addr      string                 `yaml:"addr" env:"STORAGE_SYSTEM_GRPC_ADDR" desc:"The bind address of the GRPC service."`
	TLS       *shared.GRPCServiceTLS `yaml:"tls"`
	Namespace string                 `yaml:"-"`
	Protocol  string                 `yaml:"protocol" env:"STORAGE_SYSTEM_GRPC_PROTOCOL" desc:"The transport protocol of the GPRC service."`
}

// HTTPConfig holds HTTPConfig config
type HTTPConfig struct {
	Addr      string `yaml:"addr" env:"STORAGE_SYSTEM_HTTP_ADDR" desc:"The bind address of the HTTP service."`
	Namespace string `yaml:"-"`
	Protocol  string `yaml:"protocol" env:"STORAGE_SYSTEM_HTTP_PROTOCOL" desc:"The transport protocol of the HTTP service."`
}

// Drivers holds Drivers config
type Drivers struct {
	OCIS OCISDriver `yaml:"ocis"`
}

// OCISDriver holds ocis Driver config
type OCISDriver struct {
	MetadataBackend string `yaml:"metadata_backend" env:"OCIS_DECOMPOSEDFS_METADATA_BACKEND;STORAGE_SYSTEM_OCIS_METADATA_BACKEND" desc:"The backend to use for storing metadata. Supported values are 'messagepack' and 'xattrs'. The setting 'messagepack' uses a dedicated file to store file metadata while 'xattrs' uses extended attributes to store file metadata. Defaults to 'messagepack'."`
	// Root is the absolute path to the location of the data
	Root string `yaml:"root" env:"STORAGE_SYSTEM_OCIS_ROOT" desc:"Path for the directory where the STORAGE-SYSTEM service stores it's persistent data. If not defined, the root directory derives from $OCIS_BASE_DATA_PATH:/storage."`

	MaxAcquireLockCycles    int `yaml:"max_acquire_lock_cycles" env:"STORAGE_SYSTEM_OCIS_MAX_ACQUIRE_LOCK_CYCLES" desc:"When trying to lock files, ocis will try this amount of times to acquire the lock before failing. After each try it will wait for an increasing amount of time. Values of 0 or below will be ignored and the default value of 20 will be used."`
	LockCycleDurationFactor int `yaml:"lock_cycle_duration_factor" env:"STORAGE_SYSTEM_OCIS_LOCK_CYCLE_DURATION_FACTOR" desc:"When trying to lock files, ocis will multiply the cycle with this factor and use it as a millisecond timeout. Values of 0 or below will be ignored and the default value of 30 will be used."`
}

// Cache holds cache config
type Cache struct {
	Store    string        `yaml:"store" env:"OCIS_CACHE_STORE;STORAGE_SYSTEM_CACHE_STORE" desc:"The type of the cache store. Supported values are: 'memory', 'ocmem', 'etcd', 'redis', 'redis-sentinel', 'nats-js', 'noop'. See the text description for details."`
	Nodes    []string      `yaml:"nodes" env:"OCIS_CACHE_STORE_NODES;STORAGE_SYSTEM_CACHE_STORE_NODES" desc:"A list of nodes to access the configured store. This has no effect when 'memory' or 'ocmem' stores are configured. Note that the behaviour how nodes are used is dependent on the library of the configured store. See the Environment Variable Types description for more details."`
	Database string        `yaml:"database" env:"OCIS_CACHE_DATABASE" desc:"The database name the configured store should use."`
	TTL      time.Duration `yaml:"ttl" env:"OCIS_CACHE_TTL;STORAGE_SYSTEM_CACHE_TTL" desc:"Default time to live for user info in the user info cache. Only applied when access tokens has no expiration. See the Environment Variable Types description for more details."`
	Size     int           `yaml:"size" env:"OCIS_CACHE_SIZE;STORAGE_SYSTEM_CACHE_SIZE" desc:"The maximum quantity of items in the user info cache. Only applies when store type 'ocmem' is configured. Defaults to 512 which is derived from the ocmem package though not exclicitely set as default."`
}
