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
	AuthProvider          string        `yaml:"auth_provider" env:"AUTH_MACHINE_AUTH_PROVIDER" desc:"The auth provider which should be used by the service"`
	AuthProviders         AuthProviders `yaml:"auth_providers"`
}
type Tracing struct {
	Enabled   bool   `yaml:"enabled" env:"OCIS_TRACING_ENABLED;AUTH_MACHINE_TRACING_ENABLED" desc:"Activates tracing."`
	Type      string `yaml:"type" env:"OCIS_TRACING_TYPE;AUTH_MACHINE_TRACING_TYPE"`
	Endpoint  string `yaml:"endpoint" env:"OCIS_TRACING_ENDPOINT;AUTH_MACHINE_TRACING_ENDPOINT" desc:"The endpoint to the tracing collector."`
	Collector string `yaml:"collector" env:"OCIS_TRACING_COLLECTOR;AUTH_MACHINE_TRACING_COLLECTOR"`
}

type Logging struct {
	Level  string `yaml:"level" env:"OCIS_LOG_LEVEL;AUTH_MACHINE_LOG_LEVEL" desc:"The log level."`
	Pretty bool   `yaml:"pretty" env:"OCIS_LOG_PRETTY;AUTH_MACHINE_LOG_PRETTY" desc:"Activates pretty log output."`
	Color  bool   `yaml:"color" env:"OCIS_LOG_COLOR;AUTH_MACHINE_LOG_COLOR" desc:"Activates colorized log output."`
	File   string `yaml:"file" env:"OCIS_LOG_FILE;AUTH_MACHINE_LOG_FILE" desc:"The target log file."`
}

type Service struct {
	Name string `yaml:"-"`
}

type Debug struct {
	Addr   string `yaml:"addr" env:"AUTH_MACHINE_DEBUG_ADDR"`
	Token  string `yaml:"token" env:"AUTH_MACHINE_DEBUG_TOKEN"`
	Pprof  bool   `yaml:"pprof" env:"AUTH_MACHINE_DEBUG_PPROF"`
	Zpages bool   `yaml:"zpages" env:"AUTH_MACHINE_DEBUG_ZPAGES"`
}

type GRPCConfig struct {
	Addr     string `yaml:"addr" env:"AUTH_MACHINE_GRPC_ADDR" desc:"The address of the grpc service."`
	Protocol string `yaml:"protocol" env:"AUTH_MACHINE_GRPC_PROTOCOL" desc:"The transport protocol of the grpc service."`
}

type AuthProviders struct {
	Machine MachineProvider `yaml:"machine"`
}

type MachineProvider struct {
	APIKey string `yaml:"api_key" env:"OCIS_MACHINE_AUTH_API_KEY;AUTH_MACHINE_PROVIDER_API_KEY" desc:"The api key for the machine auth provider."`
}
