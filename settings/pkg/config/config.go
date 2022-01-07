package config

import (
	"context"

	"github.com/owncloud/ocis/ocis-pkg/shared"
)

// Config combines all available configuration parts.
type Config struct {
	*shared.Commons

	Service Service

	Tracing Tracing `ocisConfig:"tracing"`
	Log     *Log    `ocisConfig:"log"`
	Debug   Debug   `ocisConfig:"debug"`

	HTTP HTTP `ocisConfig:"http"`
	GRPC GRPC `ocisConfig:"grpc"`

	DataPath     string       `ocisConfig:"data_path" env:"SETTINGS_DATA_PATH"`
	Asset        Asset        `ocisConfig:"asset"`
	TokenManager TokenManager `ocisConfig:"token_manager"`

	Context context.Context
}

// Asset defines the available asset configuration.
type Asset struct {
	Path string `ocisConfig:"path" env:"SETTINGS_ASSET_PATH"`
}
