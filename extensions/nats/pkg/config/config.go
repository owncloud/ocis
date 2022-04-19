package config

import (
	"context"

	"github.com/owncloud/ocis/ocis-pkg/shared"
)

// Config combines all available configuration parts.
type Config struct {
	*shared.Commons `yaml:"-"`

	Service Service `yaml:"-"`

	Log   *Log  `yaml:"log,omitempty"`
	Debug Debug `yaml:"debug,omitempty"`

	Nats Nats `ociConfig:"nats,omitempty"`

	Context context.Context `yaml:"-"`
}

// Nats is the nats config
type Nats struct {
	Host      string `yaml:"host,omitempty" env:"NATS_NATS_HOST"`
	Port      int    `yaml:"port,omitempty" env:"NATS_NATS_PORT"`
	ClusterID string `yaml:"clusterid,omitempty" env:"NATS_NATS_CLUSTER_ID"`
	StoreDir  string `yaml:"store_dir,omitempty" env:"NATS_NATS_STORE_DIR"`
}
