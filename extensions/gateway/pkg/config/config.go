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

	CommitShareToStorageGrant  bool
	CommitShareToStorageRef    bool
	ShareFolder                string
	DisableHomeCreationOnLogin bool
	TransferSecret             string `env:"STORAGE_TRANSFER_SECRET"`
	TransferExpires            int
	HomeMapping                string
	EtagCacheTTL               int

	UsersEndpoint             string
	GroupsEndpoint            string
	PermissionsEndpoint       string
	SharingEndpoint           string
	DataGatewayPublicURL      string
	FrontendPublicURL         string `env:"OCIS_URL;GATEWAY_FRONTEND_PUBLIC_URL"`
	AuthBasicEndpoint         string
	AuthBearerEndpoint        string
	AuthMachineEndpoint       string
	StoragePublicLinkEndpoint string
	StorageUsersEndpoint      string
	StorageSharesEndpoint     string

	StorageRegistry StorageRegistry
	AppRegistry     AppRegistry
}
type Tracing struct {
	Enabled   bool   `yaml:"enabled" env:"OCIS_TRACING_ENABLED;GATEWAY_TRACING_ENABLED" desc:"Activates tracing."`
	Type      string `yaml:"type" env:"OCIS_TRACING_TYPE;GATEWAY_TRACING_TYPE"`
	Endpoint  string `yaml:"endpoint" env:"OCIS_TRACING_ENDPOINT;GATEWAY_TRACING_ENDPOINT" desc:"The endpoint to the tracing collector."`
	Collector string `yaml:"collector" env:"OCIS_TRACING_COLLECTOR;GATEWAY_TRACING_COLLECTOR"`
}

type Logging struct {
	Level  string `yaml:"level" env:"OCIS_LOG_LEVEL;GATEWAY_LOG_LEVEL" desc:"The log level."`
	Pretty bool   `yaml:"pretty" env:"OCIS_LOG_PRETTY;GATEWAY_LOG_PRETTY" desc:"Activates pretty log output."`
	Color  bool   `yaml:"color" env:"OCIS_LOG_COLOR;GATEWAY_LOG_COLOR" desc:"Activates colorized log output."`
	File   string `yaml:"file" env:"OCIS_LOG_FILE;GATEWAY_LOG_FILE" desc:"The target log file."`
}

type Service struct {
	Name string `yaml:"-"`
}

type Debug struct {
	Addr   string `yaml:"addr" env:"GATEWAY_DEBUG_ADDR"`
	Token  string `yaml:"token" env:"GATEWAY_DEBUG_TOKEN"`
	Pprof  bool   `yaml:"pprof" env:"GATEWAY_DEBUG_PPROF"`
	Zpages bool   `yaml:"zpages" env:"GATEWAY_DEBUG_ZPAGES"`
}

type GRPCConfig struct {
	Addr     string `yaml:"addr" env:"GATEWAY_GRPC_ADDR" desc:"The address of the grpc service."`
	Protocol string `yaml:"protocol" env:"GATEWAY_GRPC_PROTOCOL" desc:"The transport protocol of the grpc service."`
}

type StorageRegistry struct {
	Driver string
	Rules  []string
	JSON   string
}

type AppRegistry struct {
	MimetypesJSON string
}
