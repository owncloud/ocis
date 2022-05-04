package config

import (
	"context"

	"github.com/owncloud/ocis/ocis-pkg/shared"
)

// Config combines all available configuration parts.
type Config struct {
	*shared.Commons `yaml:"-"`

	Service Service `yaml:"-"`

	Tracing *Tracing `yaml:"tracing"`
	Log     *Log     `yaml:"log"`
	Debug   Debug    `yaml:"debug"`

	HTTP HTTP `yaml:"http"`
	GRPC GRPC `yaml:"grpc"`

	StoreType string   `yaml:"store_type" env:"SETTINGS_STORE_TYPE"`
	DataPath  string   `yaml:"data_path" env:"SETTINGS_DATA_PATH"`
	Metadata  Metadata `yaml:"metadata_config"`

	Asset        Asset         `yaml:"asset"`
	TokenManager *TokenManager `yaml:"token_manager"`

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

	SystemUserID      string `yaml:"system_user_id" env:"OCIS_SYSTEM_USER_ID;METADATA_SYSTEM_USER_UUID"`
	SystemUserIDP     string `yaml:"system_user_idp" env:"OCIS_URL;OCIS_SYSTEM_USER_IDP;METADATA_SYSTEM_USER_IDP"`
	MachineAuthAPIKey string `yaml:"machine_auth_api_key" env:"OCIS_MACHINE_AUTH_API_KEY"`
}
