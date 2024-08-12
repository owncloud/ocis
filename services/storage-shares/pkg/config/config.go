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

	GRPC GRPCConfig `yaml:"grpc"`

	TokenManager *TokenManager `yaml:"token_manager"`
	Reva         *shared.Reva  `yaml:"reva"`

	SkipUserGroupsInToken bool `yaml:"skip_user_groups_in_token" env:"STORAGE_SHARES_SKIP_USER_GROUPS_IN_TOKEN" desc:"Disables the loading of user's group memberships from the reva access token." introductionVersion:"pre5.0"`

	MountID                string `yaml:"mount_id" env:"STORAGE_SHARES_MOUNT_ID" desc:"Mount ID of this storage. Admins can set the ID for the storage in this config option manually which is then used to reference the storage. Any reasonable long string is possible, preferably this would be an UUIDv4 format." introductionVersion:"pre5.0"`
	ReadOnly               bool   `yaml:"readonly" env:"STORAGE_SHARES_READ_ONLY" desc:"Set this storage to be read-only." introductionVersion:"pre5.0"`
	SharesProviderEndpoint string `yaml:"user_share_provider_endpoint" env:"STORAGE_SHARES_USER_SHARE_PROVIDER_ENDPOINT" desc:"GRPC endpoint of the SHARING service." introductionVersion:"pre5.0"`

	Supervised bool            `yaml:"-"`
	Context    context.Context `yaml:"-"`
}
type Log struct {
	Level  string `yaml:"level" env:"OCIS_LOG_LEVEL;STORAGE_SHARES_LOG_LEVEL" desc:"The log level. Valid values are: 'panic', 'fatal', 'error', 'warn', 'info', 'debug', 'trace'." introductionVersion:"pre5.0"`
	Pretty bool   `yaml:"pretty" env:"OCIS_LOG_PRETTY;STORAGE_SHARES_LOG_PRETTY" desc:"Activates pretty log output." introductionVersion:"pre5.0"`
	Color  bool   `yaml:"color" env:"OCIS_LOG_COLOR;STORAGE_SHARES_LOG_COLOR" desc:"Activates colorized log output." introductionVersion:"pre5.0"`
	File   string `yaml:"file" env:"OCIS_LOG_FILE;STORAGE_SHARES_LOG_FILE" desc:"The path to the log file. Activates logging to this file if set." introductionVersion:"pre5.0"`
}

type Service struct {
	Name string `yaml:"-"`
}

type Debug struct {
	Addr   string `yaml:"addr" env:"STORAGE_SHARES_DEBUG_ADDR" desc:"Bind address of the debug server, where metrics, health, config and debug endpoints will be exposed." introductionVersion:"pre5.0"`
	Token  string `yaml:"token" env:"STORAGE_SHARES_DEBUG_TOKEN" desc:"Token to secure the metrics endpoint." introductionVersion:"pre5.0"`
	Pprof  bool   `yaml:"pprof" env:"STORAGE_SHARES_DEBUG_PPROF" desc:"Enables pprof, which can be used for profiling." introductionVersion:"pre5.0"`
	Zpages bool   `yaml:"zpages" env:"STORAGE_SHARES_DEBUG_ZPAGES" desc:"Enables zpages, which can be used for collecting and viewing in-memory traces." introductionVersion:"pre5.0"`
}

type GRPCConfig struct {
	Addr      string                 `yaml:"addr" env:"STORAGE_SHARES_GRPC_ADDR" desc:"The bind address of the GRPC service." introductionVersion:"pre5.0"`
	TLS       *shared.GRPCServiceTLS `yaml:"tls"`
	Namespace string                 `yaml:"-"`
	Protocol  string                 `yaml:"protocol" env:"OCIS_GRPC_PROTOCOL;STORAGE_SHARES_GRPC_PROTOCOL" desc:"The transport protocol of the GRPC service." introductionVersion:"pre5.0"`
}
