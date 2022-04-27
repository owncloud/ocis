package config

import "github.com/owncloud/ocis/ocis-pkg/shared"

type Config struct {
	*shared.Commons `yaml:"-"`

	Service    Service  `yaml:"-"`
	Tracing    *Tracing `yaml:"tracing,omitempty"`
	Logging    *Logging `yaml:"log,omitempty"`
	Debug      Debug    `yaml:"debug,omitempty"`
	Supervised bool     `yaml:"supervised,omitempty"`

	GRPC GRPCConfig `yaml:"grpc,omitempty"`

	TokenManager *TokenManager `yaml:"token_manager,omitempty"`
	Reva         *Reva         `yaml:"reva,omitempty"`

	SkipUserGroupsInToken bool `yaml:",omitempty"`

	CommitShareToStorageGrant  bool   `yaml:"commit_share_to_storage_grant,omitempty"`
	CommitShareToStorageRef    bool   `yaml:"commit_share_to_storage_ref,omitempty"`
	ShareFolder                string `yaml:"share_folder,omitempty"`
	DisableHomeCreationOnLogin bool   `yaml:"disable_home_creation_on_login,omitempty"`
	TransferSecret             string `yaml:"transfer_secret,omitempty" env:"STORAGE_TRANSFER_SECRET"`
	TransferExpires            int    `yaml:"transfer_expires,omitempty"`
	HomeMapping                string `yaml:"home_mapping,omitempty"`
	EtagCacheTTL               int    `yaml:"etag_cache_ttl,omitempty"`

	UsersEndpoint             string `yaml:"users_endpoint,omitempty"`
	GroupsEndpoint            string `yaml:"groups_endpoint,omitempty"`
	PermissionsEndpoint       string `yaml:"permissions_endpoint,omitempty"`
	SharingEndpoint           string `yaml:"sharing_endpoint,omitempty"`
	FrontendPublicURL         string `yaml:"frontend_public_url,omitempty" env:"OCIS_URL;GATEWAY_FRONTEND_PUBLIC_URL"`
	AuthBasicEndpoint         string `yaml:"auth_basic_endpoint,omitempty"`
	AuthBearerEndpoint        string `yaml:"auth_bearer_endpoint,omitempty"`
	AuthMachineEndpoint       string `yaml:"auth_machine_endpoint,omitempty"`
	StoragePublicLinkEndpoint string `yaml:"storage_public_link_endpoint,omitempty"`
	StorageUsersEndpoint      string `yaml:"storage_users_endpoint,omitempty"`
	StorageSharesEndpoint     string `yaml:"storage_shares_endpoint,omitempty"`

	StorageRegistry StorageRegistry `yaml:"storage_registry,omitempty"`
	AppRegistry     AppRegistry     `yaml:"app_registry,omitempty"`
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
