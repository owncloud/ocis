package config

import (
	"context"

	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
)

// Config combines all available configuration parts.
type Config struct {
	Commons *shared.Commons `yaml:"-"` // don't use this directly as configuration for a service

	Service Service `yaml:"-"`

	Log   *Log  `yaml:"log"`
	Debug Debug `yaml:"debug"`

	Nats Nats `ociConfig:"nats"`

	Context context.Context `yaml:"-"`
}

// Nats is the nats config
type Nats struct {
	Host      string `yaml:"host" env:"NATS_NATS_HOST" desc:"Bind address."`
	Port      int    `yaml:"port" env:"NATS_NATS_PORT" desc:"Bind port."`
	ClusterID string `yaml:"clusterid" env:"NATS_NATS_CLUSTER_ID" desc:"ID of the NATS cluster."`
	StoreDir  string `yaml:"store_dir" env:"NATS_NATS_STORE_DIR" desc:"Path for the NATS JetStream persistence directory."`
}
