package config

import (
	"context"

	"github.com/owncloud/ocis/ocis-pkg/shared"
)

// Config combines all available configuration parts.
type Config struct {
	*shared.Commons `ocisConfig:"-" yaml:"-"`

	Service Service `ocisConfig:"-" yaml:"-"`

	Log   *Log  `ocisConfig:"log"`
	Debug Debug `ocisConfig:"debug"`

	Nats Nats `ociConfig:"nats"`

	Context context.Context `ocisConfig:"-" yaml:"-"`
}

// Nats is the nats config
type Nats struct {
	Host     string `ocisConfig:"host" env:"NATS_NATS_HOST"`
	Port     int    `ocisConfig:"port" env:"NATS_NATS_PORT"`
	StoreDir string `ocisConfig:"store_dir" env:"NATS_NATS_STORE_DIR"`
}
