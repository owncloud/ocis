package config

import (
	"context"

	"github.com/owncloud/ocis/ocis-pkg/shared"
)

// Config combines all available configuration parts.
type Config struct {
	*shared.Commons `yaml:"-"`

	Service Service `yaml:"-"`

	Tracing *Tracing `yaml:"tracing,omitempty"`
	Log     *Log     `yaml:"log,omitempty"`
	Debug   Debug    `yaml:"debug,omitempty"`

	HTTP HTTP `yaml:"http,omitempty"`
	GRPC GRPC `yaml:"grpc,omitempty"`

	StoreType string   `yaml:"store_type,omitempty" env:"SETTINGS_STORE_TYPE"`
	DataPath  string   `yaml:"data_path,omitempty" env:"SETTINGS_DATA_PATH"`
	Metadata  Metadata `yaml:"metadata_config,omitempty"`

	Asset        Asset        `yaml:"asset,omitempty"`
	TokenManager TokenManager `yaml:"token_manager,omitempty"`

	Context context.Context `yaml:"-"`
}

// Asset defines the available asset configuration.
type Asset struct {
	Path string `yaml:"path" env:"SETTINGS_ASSET_PATH"`
}

// Metadata configures the metadata store to use
type Metadata struct {
	GatewayAddress string `yaml:"gateway_addr" env:"STORAGE_GATEWAY_GRPC_ADDR"`
	StorageAddress string `yaml:"storage_addr" env:"STORAGE_GRPC_ADDR"`

	ServiceUserID     string `yaml:"service_user_id" env:"METADATA_SERVICE_USER_UUID"`
	ServiceUserIDP    string `yaml:"service_user_idp" env:"OCIS_URL;METADATA_SERVICE_USER_IDP"`
	MachineAuthAPIKey string `yaml:"machine_auth_api_key" env:"OCIS_MACHINE_AUTH_API_KEY"`
}
