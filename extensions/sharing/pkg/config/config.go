package config

import "github.com/owncloud/ocis/ocis-pkg/shared"

type Config struct {
	*shared.Commons `yaml:"-"`
	Service         Service  `yaml:"-"`
	Tracing         *Tracing `yaml:"tracing"`
	Logging         *Logging `yaml:"log"`
	Debug           Debug    `yaml:"debug"`
	Supervised      bool

	GRPC GRPCConfig `yaml:"grpc"`

	JWTSecret             string
	GatewayEndpoint       string
	SkipUserGroupsInToken bool
	UserSharingDriver     string
	UserSharingDrivers    UserSharingDrivers
	PublicSharingDriver   string
	PublicSharingDrivers  PublicSharingDrivers
	Events                Events
}
type Tracing struct {
	Enabled   bool   `yaml:"enabled" env:"OCIS_TRACING_ENABLED;SHARING_TRACING_ENABLED" desc:"Activates tracing."`
	Type      string `yaml:"type" env:"OCIS_TRACING_TYPE;SHARING_TRACING_TYPE"`
	Endpoint  string `yaml:"endpoint" env:"OCIS_TRACING_ENDPOINT;SHARING_TRACING_ENDPOINT" desc:"The endpoint to the tracing collector."`
	Collector string `yaml:"collector" env:"OCIS_TRACING_COLLECTOR;SHARING_TRACING_COLLECTOR"`
}

type Logging struct {
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
	Addr     string `yaml:"addr" env:"SHARING_GRPC_ADDR" desc:"The address of the grpc service."`
	Protocol string `yaml:"protocol" env:"SHARING_GRPC_PROTOCOL" desc:"The transport protocol of the grpc service."`
}

type UserSharingDrivers struct {
	JSON UserSharingJSONDriver
	SQL  UserSharingSQLDriver
	CS3  UserSharingCS3Driver
}

type UserSharingJSONDriver struct {
	File string
}

type UserSharingSQLDriver struct {
	DBUsername                 string
	DBPassword                 string
	DBHost                     string
	DBPort                     int
	DBName                     string
	PasswordHashCost           int
	EnableExpiredSharesCleanup bool
	JanitorRunInterval         int
	UserStorageMountID         string
}

type UserSharingCS3Driver struct {
	ProviderAddr      string
	ServiceUserID     string
	ServiceUserIDP    string
	MachineAuthAPIKey string
}

type PublicSharingDrivers struct {
	JSON PublicSharingJSONDriver
	SQL  PublicSharingSQLDriver
	CS3  PublicSharingCS3Driver
}

type PublicSharingJSONDriver struct {
	File string
}

type PublicSharingSQLDriver struct {
	DBUsername                 string
	DBPassword                 string
	DBHost                     string
	DBPort                     int
	DBName                     string
	PasswordHashCost           int
	EnableExpiredSharesCleanup bool
	JanitorRunInterval         int
	UserStorageMountID         string
}

type PublicSharingCS3Driver struct {
	ProviderAddr      string
	ServiceUserID     string
	ServiceUserIDP    string
	MachineAuthAPIKey string
}

type Events struct {
	Addr      string
	ClusterID string
}
