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

	GRPC GRPC `ocisConfig:"grpc"`

	Datapath string `ocisConfig:"data_path" env:"STORE_DATA_PATH"`

	Context context.Context `ocisConfig:"-" yaml:"-"`
}
