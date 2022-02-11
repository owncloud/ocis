package config

import (
	"context"

	"github.com/owncloud/ocis/ocis-pkg/shared"
)

// Config combines all available configuration parts.
type Config struct {
	*shared.Commons

	Service Service

	Log   *Log  `ocisConfig:"log"`
	Debug Debug `ocisConfig:"debug"`

	Nats Nats `ociConfig:"nats"`

	Context context.Context
}

// Nats is the nats config
type Nats struct {
	Host string
	Port int
}
