package config

import (
	"context"

	"github.com/owncloud/ocis/ocis-pkg/shared"
)

// Config combines all available configuration parts.
type Config struct {
	*shared.Commons

	Service Service

	Tracing *Tracing `ocisConfig:"tracing"`
	Log     *Log     `ocisConfig:"log"`
	Debug   Debug    `ocisConfig:"debug"`

	Nats Nats `ociConfig:"nats"`

	HTTP HTTP `ocisConfig:"http"`
	GRPC GRPC `ocisConfig:"grpc"`

	DataPath     string       `ocisConfig:"data_path" env:"SETTINGS_DATA_PATH"`
	TokenManager TokenManager `ocisConfig:"token_manager"`

	Context context.Context
}

// Nats is the nats config
type Nats struct {
	Host string
	Port int
}
