package config

import (
	"context"
	"time"

	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
)

type Config struct {
	Commons *shared.Commons `yaml:"-"` // don't use this directly as configuration for a service
	Service Service         `yaml:"-"`
	Tracing *Tracing        `yaml:"tracing"`
	Log     *Log            `yaml:"log"`
	Debug   Debug           `yaml:"debug"`

	GRPC GRPCConfig `yaml:"grpc"`

	TokenManager *TokenManager `yaml:"token_manager"`
	Reva         *shared.Reva  `yaml:"reva"`
	BruteForce   BruteForce    `yaml:"brute_force"`
	Store        Store         `yaml:"store"`

	SkipUserGroupsInToken bool `yaml:"skip_user_groups_in_token" env:"STORAGE_PUBLICLINK_SKIP_USER_GROUPS_IN_TOKEN" desc:"Disables the loading of user's group memberships from the reva access token." introductionVersion:"pre5.0"`

	StorageProvider StorageProvider `yaml:"storage_provider"`

	Context context.Context `yaml:"-"`
}

type Log struct {
	Level  string `yaml:"level" env:"OCIS_LOG_LEVEL;STORAGE_PUBLICLINK_LOG_LEVEL" desc:"The log level. Valid values are: 'panic', 'fatal', 'error', 'warn', 'info', 'debug', 'trace'." introductionVersion:"pre5.0"`
	Pretty bool   `yaml:"pretty" env:"OCIS_LOG_PRETTY;STORAGE_PUBLICLINK_LOG_PRETTY" desc:"Activates pretty log output." introductionVersion:"pre5.0"`
	Color  bool   `yaml:"color" env:"OCIS_LOG_COLOR;STORAGE_PUBLICLINK_LOG_COLOR" desc:"Activates colorized log output." introductionVersion:"pre5.0"`
	File   string `yaml:"file" env:"OCIS_LOG_FILE;STORAGE_PUBLICLINK_LOG_FILE" desc:"The path to the log file. Activates logging to this file if set." introductionVersion:"pre5.0"`
}

type Service struct {
	Name string `yaml:"-"`
}

type Debug struct {
	Addr   string `yaml:"addr" env:"STORAGE_PUBLICLINK_DEBUG_ADDR" desc:"Bind address of the debug server, where metrics, health, config and debug endpoints will be exposed." introductionVersion:"pre5.0"`
	Token  string `yaml:"token" env:"STORAGE_PUBLICLINK_DEBUG_TOKEN" desc:"Token to secure the metrics endpoint." introductionVersion:"pre5.0"`
	Pprof  bool   `yaml:"pprof" env:"STORAGE_PUBLICLINK_DEBUG_PPROF" desc:"Enables pprof, which can be used for profiling." introductionVersion:"pre5.0"`
	Zpages bool   `yaml:"zpages" env:"STORAGE_PUBLICLINK_DEBUG_ZPAGES" desc:"Enables zpages, which can be used for collecting and viewing in-memory traces." introductionVersion:"pre5.0"`
}

type GRPCConfig struct {
	Addr      string                 `yaml:"addr" env:"STORAGE_PUBLICLINK_GRPC_ADDR" desc:"The bind address of the GRPC service." introductionVersion:"pre5.0"`
	TLS       *shared.GRPCServiceTLS `yaml:"tls"`
	Namespace string                 `yaml:"-"`
	Protocol  string                 `yaml:"protocol" env:"OCIS_GRPC_PROTOCOL;STORAGE_PUBLICLINK_GRPC_PROTOCOL" desc:"The transport protocol of the GRPC service." introductionVersion:"pre5.0"`
}

type StorageProvider struct {
	MountID string `yaml:"mount_id" env:"STORAGE_PUBLICLINK_STORAGE_PROVIDER_MOUNT_ID" desc:"Mount ID of this storage. Admins can set the ID for the storage in this config option manually which is then used to reference the storage. Any reasonable long string is possible, preferably this would be an UUIDv4 format." introductionVersion:"pre5.0"`
}

type BruteForce struct {
	TimeGap     time.Duration `yaml:"time_gap" env:"STORAGE_PUBLICLINK_BRUTEFORCE_TIMEGAP" desc:"The duration of the time gap computed for the brute force protection." introductionVersion:"Curie"`
	MaxAttempts int           `yaml:"max_attempts" env:"STORAGE_PUBLICLINK_BRUTEFORCE_MAXATTEMPTS" desc:"The maximum number of failed attempts allowed in the time gap defined in STORAGE_PUBLICLINK_BRUTEFORCE_TIMEGAP." introductionVersion:"Curie"`
}

// Store configures the store to use
type Store struct {
	Store        string   `yaml:"store" env:"OCIS_PERSISTENT_STORE;STORAGE_PUBLICLINK_STORE_STORE" desc:"The type of the store. Supported values are: 'memory', 'nats-js-kv', 'redis-sentinel', 'noop'. See the text description for details." introductionVersion:"curie"`
	Nodes        []string `yaml:"nodes" env:"OCIS_PERSISTENT_STORE_NODES;STORAGE_PUBLICLINK_STORE_NODES" desc:"A list of nodes to access the configured store. This has no effect when 'memory' store is configured. Note that the behaviour how nodes are used is dependent on the library of the configured store. See the Environment Variable Types description for more details." introductionVersion:"curie"`
	Database     string   `yaml:"database" env:"STORAGE_PUBLICLINK_STORE_DATABASE" desc:"The database name the configured store should use." introductionVersion:"curie"`
	Table        string   `yaml:"table" env:"STORAGE_PUBLICLINK_STORE_TABLE" desc:"The database table the store should use." introductionVersion:"curie"`
	AuthUsername string   `yaml:"username" env:"OCIS_PERSISTENT_STORE_AUTH_USERNAME;STORAGE_PUBLICLINK_STORE_AUTH_USERNAME" desc:"The username to authenticate with the store. Only applies when store type 'nats-js-kv' is configured." introductionVersion:"curie"`
	AuthPassword string   `yaml:"password" env:"OCIS_PERSISTENT_STORE_AUTH_PASSWORD;STORAGE_PUBLICLINK_STORE_AUTH_PASSWORD" desc:"The password to authenticate with the store. Only applies when store type 'nats-js-kv' is configured." introductionVersion:"curie"`
}
