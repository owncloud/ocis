package config

import (
	"context"
)

// Config combines all available configuration parts.
type Config struct {
	Service Service

	Tracing Tracing `ocisConfig:"tracing"`
	Log     Log     `ocisConfig:"log"`
	Debug   Debug   `ocisConfig:"debug"`

	GRPC GRPC `ocisConfig:"grpc"`

	Datapath string `ocisConfig:"data_path" env:"STORE_DATA_PATH"`

	Context    context.Context
	Supervised bool
}
