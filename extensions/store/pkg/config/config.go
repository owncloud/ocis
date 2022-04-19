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

	GRPC GRPC `yaml:"grpc"`

	Datapath string `yaml:"data_path" env:"STORE_DATA_PATH"`

	ConfigFile string `yaml:"-" env:"STORE_CONFIG_FILE" desc:"config file to be used by the store extension"`

	Context context.Context `yaml:"-"`
}
