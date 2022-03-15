package config

import (
	"context"

	"github.com/owncloud/ocis/ocis-pkg/shared"
)

// Config combines all available configuration parts.
type Config struct {
	*shared.Commons `ocisConfig:"-" yaml:"-"`

	Service Service `ocisConfig:"-" yaml:"-"`

	Tracing *Tracing `ocisConfig:"tracing"`
	Log     *Log     `ocisConfig:"log"`
	Debug   Debug    `ocisConfig:"debug"`

	HTTP HTTP `ocisConfig:"http"`
	GRPC GRPC `ocisConfig:"grpc"`

	StoreType string   `ocisConfig:"store_type" env:"SETTINGS_STORE_TYPE"`
	DataPath  string   `ocisConfig:"data_path" env:"SETTINGS_DATA_PATH"`
	Metadata  Metadata `ocisConfig:"metadata_config"`

	Asset        Asset        `ocisConfig:"asset"`
	TokenManager TokenManager `ocisConfig:"token_manager"`

	Context context.Context `ocisConfig:"-" yaml:"-"`
}

// Asset defines the available asset configuration.
type Asset struct {
	Path string `ocisConfig:"path" env:"SETTINGS_ASSET_PATH"`
}

// Metadata configures the metadata store to use
type Metadata struct {
	GatewayAddress string `ocisConfig:"gateway_addr" env:"STORAGE_GATEWAY_GRPC_ADDR"`
	StorageAddress string `ocisConfig:"storage_addr" env:"STORAGE_GRPC_ADDR"`

	ServiceUserID     string `ocisConfig:"service_user_id" env:"METADATA_SERVICE_USER_UUID"`
	ServiceUserIDP    string `ocisConfig:"service_user_idp" env:"OCIS_URL;METADATA_SERVICE_USER_IDP"`
	MachineAuthAPIKey string `ocisConfig:"machine_auth_api_key" env:"OCIS_MACHINE_AUTH_API_KEY"`
}
