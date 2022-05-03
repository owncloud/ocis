package config

import (
	"context"

	"github.com/owncloud/ocis/ocis-pkg/shared"
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
	Events       Events        `yaml:"events"`

	SkipUserGroupsInToken bool `yaml:"-"`

	UserSharingDriver    string               `yaml:"user_sharing_driver" env:"SHARING_USER_DRIVER"`
	UserSharingDrivers   UserSharingDrivers   `yaml:"user_sharin_drivers"`
	PublicSharingDriver  string               `yaml:"public_sharing_driver" env:"SHARING_PUBLIC_DRIVER"`
	PublicSharingDrivers PublicSharingDrivers `yaml:"public_sharing_drivers"`

	Supervised bool            `yaml:"-"`
	Context    context.Context `yaml:"-"`
}
type Tracing struct {
	Enabled   bool   `yaml:"enabled" env:"OCIS_TRACING_ENABLED;SHARING_TRACING_ENABLED" desc:"Activates tracing."`
	Type      string `yaml:"type" env:"OCIS_TRACING_TYPE;SHARING_TRACING_TYPE"`
	Endpoint  string `yaml:"endpoint" env:"OCIS_TRACING_ENDPOINT;SHARING_TRACING_ENDPOINT" desc:"The endpoint to the tracing collector."`
	Collector string `yaml:"collector" env:"OCIS_TRACING_COLLECTOR;SHARING_TRACING_COLLECTOR"`
}

type Log struct {
	Level  string `yaml:"level" env:"OCIS_LOG_LEVEL;SHARING_LOG_LEVEL" desc:"The log level."`
	Pretty bool   `yaml:"pretty" env:"OCIS_LOG_PRETTY;SHARING_LOG_PRETTY" desc:"Activates pretty log output."`
	Color  bool   `yaml:"color" env:"OCIS_LOG_COLOR;SHARING_LOG_COLOR" desc:"Activates colorized log output."`
	File   string `yaml:"file" env:"OCIS_LOG_FILE;SHARING_LOG_FILE" desc:"The target log file."`
}

type Service struct {
	Name string `yaml:"-"`
}

type Debug struct {
	Addr   string `yaml:"addr" env:"SHARING_DEBUG_ADDR"`
	Token  string `yaml:"token" env:"SHARING_DEBUG_TOKEN"`
	Pprof  bool   `yaml:"pprof" env:"SHARING_DEBUG_PPROF"`
	Zpages bool   `yaml:"zpages" env:"SHARING_DEBUG_ZPAGES"`
}

type GRPCConfig struct {
	Addr      string `yaml:"addr" env:"SHARING_GRPC_ADDR" desc:"The address of the grpc service."`
	Namespace string `yaml:"-"`
	Protocol  string `yaml:"protocol" env:"SHARING_GRPC_PROTOCOL" desc:"The transport protocol of the grpc service."`
}

type UserSharingDrivers struct {
	JSON UserSharingJSONDriver `yaml:"json"`
	CS3  UserSharingCS3Driver  `yaml:"cs3"`
	SQL  UserSharingSQLDriver  `yaml:"sql,omitempty"` // not supported by the oCIS product, therefore not part of docs
}

type UserSharingJSONDriver struct {
	File string `yaml:"file" env:"SHARING_USER_JSON_FILE"`
}

type UserSharingSQLDriver struct {
	DBUsername                 string `yaml:"db_username"`
	DBPassword                 string `yaml:"db_password"`
	DBHost                     string `yaml:"db_host"`
	DBPort                     int    `yaml:"db_port"`
	DBName                     string `yaml:"db_name"`
	PasswordHashCost           int    `yaml:"password_hash_cost"`
	EnableExpiredSharesCleanup bool   `yaml:"enable_expired_shares_cleanup"`
	JanitorRunInterval         int    `yaml:"janitor_run_interval"`
	UserStorageMountID         string `yaml:"user_storage_mount_id"`
}

type UserSharingCS3Driver struct {
	ProviderAddr      string `yaml:"provider_addr" env:"SHARING_USER_CS3_PROVIDER_ADDR"`
	ServiceUserID     string `yaml:"service_user_id" env:"SHARING_USER_CS3_SERVICE_USER_ID"`
	ServiceUserIDP    string `yaml:"service_user_idp" env:"OCIS_URL;SHARING_USER_CS3_SERVICE_USER_IDP"`
	MachineAuthAPIKey string `yaml:"machine_auth_api_key" env:"OCIS_MACHINE_AUTH_API_KEY;SHARING_USER_CS3_MACHINE_AUTH_API_KEY"`
}

type PublicSharingDrivers struct {
	JSON PublicSharingJSONDriver `yaml:"json"`
	CS3  PublicSharingCS3Driver  `yaml:"cs3"`
	SQL  PublicSharingSQLDriver  `yaml:"sql,omitempty"` // not supported by the oCIS product, therefore not part of docs
}

type PublicSharingJSONDriver struct {
	File string `yaml:"file" env:"SHARING_PUBLIC_JSON_FILE"`
}

type PublicSharingSQLDriver struct {
	DBUsername                 string `yaml:"db_username"`
	DBPassword                 string `yaml:"db_password"`
	DBHost                     string `yaml:"db_host"`
	DBPort                     int    `yaml:"db_port"`
	DBName                     string `yaml:"db_name"`
	PasswordHashCost           int    `yaml:"password_hash_cost"`
	EnableExpiredSharesCleanup bool   `yaml:"enable_expired_shares_cleanup"`
	JanitorRunInterval         int    `yaml:"janitor_run_interval"`
	UserStorageMountID         string `yaml:"user_storage_mount_id"`
}

type PublicSharingCS3Driver struct {
	ProviderAddr      string `yaml:"provider_addr" env:"SHARING_PUBLIC_CS3_PROVIDER_ADDR"`
	ServiceUserID     string `yaml:"service_user_id" env:"SHARING_PUBLIC_CS3_SERVICE_USER_ID"`
	ServiceUserIDP    string `yaml:"service_user_idp" env:"OCIS_URL;SHARING_PUBLIC_CS3_SERVICE_USER_IDP"`
	MachineAuthAPIKey string `yaml:"machine_auth_api_key" env:"OCIS_MACHINE_AUTH_API_KEY;SHARING_PUBLIC_CS3_MACHINE_AUTH_API_KEY"`
}

type Events struct {
	Addr      string `yaml:"addr" env:"SHARING_EVENTS_ADDR"`
	ClusterID string `yaml:"cluster_id" env:"SHARING_EVENTS_CLUSTER_ID"`
}
