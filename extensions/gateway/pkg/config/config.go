package config

import (
	"context"

	"github.com/owncloud/ocis/ocis-pkg/shared"
)

type Config struct {
	*shared.Commons `yaml:"-"`

	Service Service  `yaml:"-"`
	Tracing *Tracing `yaml:"tracing"`
	Log     *Log     `yaml:"log"`
	Debug   Debug    `yaml:"debug"`

	GRPC GRPCConfig `yaml:"grpc"`

	TokenManager *TokenManager `yaml:"token_manager"`
	Reva         *Reva         `yaml:"reva"`

	SkipUserGroupsInToken bool `yaml:"skip_user_groups_in_token" env:"GATEWAY_SKIP_USER_GROUPS_IN_TOKEN"`

	CommitShareToStorageGrant  bool   `yaml:"commit_share_to_storage_grant" env:"GATEWAY_COMMIT_SHARE_TO_STORAGE_GRANT"`
	CommitShareToStorageRef    bool   `yaml:"commit_share_to_storage_ref" env:"GATEWAY_COMMIT_SHARE_TO_STORAGE_REF"`
	ShareFolder                string `yaml:"share_folder_name" env:"GATEWAY_SHARE_FOLDER_NAME"`
	DisableHomeCreationOnLogin bool   `yaml:"disable_home_creation_on_login" env:"GATEWAY_DISABLE_HOME_CREATION_ON_LOGIN"`
	TransferSecret             string `yaml:"transfer_secret" env:"STORAGE_TRANSFER_SECRET"` // TODO: how to name the env
	TransferExpires            int    `yaml:"transfer_expires" env:"GATEWAY_TRANSFER_EXPIRES"`
	HomeMapping                string `yaml:"home_mapping" env:"GATEWAY_HOME_MAPPIN"`
	EtagCacheTTL               int    `yaml:"etag_cache_ttl" env:"GATEWAY_ETAG_CACHE_TTL"`

	FrontendPublicURL string `yaml:"frontend_public_url" env:"OCIS_URL;GATEWAY_FRONTEND_PUBLIC_URL"`

	UsersEndpoint             string `yaml:"users_endpoint" env:"GATEWAY_USERS_ENDPOINT"`
	GroupsEndpoint            string `yaml:"groups_endpoint" env:"GATEWAY_GROUPS_ENDPOINT"`
	PermissionsEndpoint       string `yaml:"permissions_endpoint" env:"GATEWAY_PERMISSIONS_ENDPOINT"`
	SharingEndpoint           string `yaml:"sharing_endpoint" env:"GATEWAY_SHARING_ENDPOINT"`
	AuthBasicEndpoint         string `yaml:"auth_basic_endpoint" env:"GATEWAY_AUTH_BASIC_ENDPOINT"`
	AuthBearerEndpoint        string `yaml:"auth_bearer_endpoint" env:"GATEWAY_AUTH_BEARER_ENDPOINT"`
	AuthMachineEndpoint       string `yaml:"auth_machine_endpoint" env:"GATEWAY_AUTH_MACHINE_ENDPOINT"`
	StoragePublicLinkEndpoint string `yaml:"storage_public_link_endpoint" env:"GATEWAY_STORAGE_PUBLIC_LINK_ENDPOINT"`
	StorageUsersEndpoint      string `yaml:"storage_users_endpoint" env:"GATEWAY_STORAGE_USERS_ENDPOINT"`
	StorageSharesEndpoint     string `yaml:"storage_shares_endpoint" env:"GATEWAY_STORAGE_SHARES_ENDPOINT"`
	AppRegistryEndpoint       string `yaml:"app_registry_endpoint" env:"GATEWAY_APP_REGISTRY_ENDPOINT"`

	StorageRegistry StorageRegistry `yaml:"storage_registry"` //TODO: should we even support switching this?

	Supervised bool            `yaml:"-"`
	Context    context.Context `yaml:"-"`
}
type Tracing struct {
	Enabled   bool   `yaml:"enabled" env:"OCIS_TRACING_ENABLED;GATEWAY_TRACING_ENABLED" desc:"Activates tracing."`
	Type      string `yaml:"type" env:"OCIS_TRACING_TYPE;GATEWAY_TRACING_TYPE"`
	Endpoint  string `yaml:"endpoint" env:"OCIS_TRACING_ENDPOINT;GATEWAY_TRACING_ENDPOINT" desc:"The endpoint to the tracing collector."`
	Collector string `yaml:"collector" env:"OCIS_TRACING_COLLECTOR;GATEWAY_TRACING_COLLECTOR"`
}

type Log struct {
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
	Addr      string `yaml:"addr" env:"GATEWAY_GRPC_ADDR" desc:"The address of the grpc service."`
	Namespace string `yaml:"-"`
	Protocol  string `yaml:"protocol" env:"GATEWAY_GRPC_PROTOCOL" desc:"The transport protocol of the grpc service."`
}

type StorageRegistry struct {
	Driver string   `yaml:"driver"` //TODO: configure via env?
	Rules  []string `yaml:"rules"`  //TODO: configure via env?
	JSON   string   `yaml:"json"`   //TODO: configure via env?
}
