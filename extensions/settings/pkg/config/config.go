package config

import (
	"context"

	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
)

// Config combines all available configuration parts.
type Config struct {
	Commons *shared.Commons `yaml:"-"` // don't use this directly as configuration for a service

	Service Service `yaml:"-"`

	Tracing *Tracing `yaml:"tracing"`
	Log     *Log     `yaml:"log"`
	Debug   Debug    `yaml:"debug"`

	HTTP HTTP `yaml:"http"`
	GRPC GRPC `yaml:"grpc"`

	StoreType string   `yaml:"store_type" env:"SETTINGS_STORE_TYPE"`
	DataPath  string   `yaml:"data_path" env:"SETTINGS_DATA_PATH"`
	Metadata  Metadata `yaml:"metadata_config"`

	AdminUserID string `yaml:"admin_user_id" env:"OCIS_ADMIN_USER_ID;SETTINGS_ADMIN_USER_ID"`

	Asset        Asset         `yaml:"asset"`
	TokenManager *TokenManager `yaml:"token_manager"`

	SetupDefaultAssignments bool `yaml:"set_default_assignments" env:"SETTINGS_SETUP_DEFAULT_ASSIGNMENTS;ACCOUNTS_DEMO_USERS_AND_GROUPS" desc:"If the default role assignments for the demo users should be setup."`

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

	SystemUserID     string `yaml:"system_user_id" env:"OCIS_SYSTEM_USER_ID;SETTINGS_SYSTEM_USER_ID"`
	SystemUserIDP    string `yaml:"system_user_idp" env:"OCIS_SYSTEM_USER_IDP;SETTINGS_SYSTEM_USER_IDP"`
	SystemUserAPIKey string `yaml:"system_user_api_key" env:"OCIS_SYSTEM_USER_API_KEY"`
}
