package config

import (
	"context"

	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
)

type Config struct {
	*shared.Commons `yaml:"-"`
	Service         Service  `yaml:"-"`
	Tracing         *Tracing `yaml:"tracing"`
	Log             *Log     `yaml:"log"`
	Debug           Debug    `yaml:"debug"`

	GRPC GRPCConfig `yaml:"grpc"`

	TokenManager *TokenManager `yaml:"token_manager"`
	Reva         *Reva         `yaml:"reva"`

	SkipUserGroupsInToken bool `yaml:"skip_user_groups_in_token" env:"STORAGE_PUBLICLINK_SKIP_USER_GROUPS_IN_TOKEN"`

	StorageProvider StorageProvider `yaml:"storage_provider"`

	Supervised bool            `yaml:"-"`
	Context    context.Context `yaml:"-"`
}
type Tracing struct {
	Enabled   bool   `yaml:"enabled" env:"OCIS_TRACING_ENABLED;STORAGE_PUBLICLINK_TRACING_ENABLED" desc:"Activates tracing."`
	Type      string `yaml:"type" env:"OCIS_TRACING_TYPE;STORAGE_PUBLICLINK_TRACING_TYPE"`
	Endpoint  string `yaml:"endpoint" env:"OCIS_TRACING_ENDPOINT;STORAGE_PUBLICLINK_TRACING_ENDPOINT" desc:"The endpoint to the tracing collector."`
	Collector string `yaml:"collector" env:"OCIS_TRACING_COLLECTOR;STORAGE_PUBLICLINK_TRACING_COLLECTOR"`
}

type Log struct {
	Level  string `yaml:"level" env:"OCIS_LOG_LEVEL;STORAGE_PUBLICLINK_LOG_LEVEL" desc:"The log level."`
	Pretty bool   `yaml:"pretty" env:"OCIS_LOG_PRETTY;STORAGE_PUBLICLINK_LOG_PRETTY" desc:"Activates pretty log output."`
	Color  bool   `yaml:"color" env:"OCIS_LOG_COLOR;STORAGE_PUBLICLINK_LOG_COLOR" desc:"Activates colorized log output."`
	File   string `yaml:"file" env:"OCIS_LOG_FILE;STORAGE_PUBLICLINK_LOG_FILE" desc:"The target log file."`
}

type Service struct {
	Name string `yaml:"-"`
}

type Debug struct {
	Addr   string `yaml:"addr" env:"STORAGE_PUBLICLINK_DEBUG_ADDR"`
	Token  string `yaml:"token" env:"STORAGE_PUBLICLINK_DEBUG_TOKEN"`
	Pprof  bool   `yaml:"pprof" env:"STORAGE_PUBLICLINK_DEBUG_PPROF"`
	Zpages bool   `yaml:"zpages" env:"STORAGE_PUBLICLINK_DEBUG_ZPAGES"`
}

type GRPCConfig struct {
	Addr      string `yaml:"addr" env:"STORAGE_PUBLICLINK_GRPC_ADDR" desc:"The address of the grpc service."`
	Namespace string `yaml:"-"`
	Protocol  string `yaml:"protocol" env:"STORAGE_PUBLICLINK_GRPC_PROTOCOL" desc:"The transport protocol of the grpc service."`
}

type StorageProvider struct {
	MountID string `yaml:"mount_id" env:"STORAGE_PUBLICLINK_STORAGE_PROVIDER_MOUNT_ID"`
}
