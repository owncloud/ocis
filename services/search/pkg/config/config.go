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

	GRPC GRPC `yaml:"grpc"`

	Reva      *Reva     `yaml:"reva"`
	Events    Events    `yaml:"events"`
	Engine    Engine    `yaml:"engine"`
	Extractor Extractor `yaml:"extractor"`

	MachineAuthAPIKey string `yaml:"machine_auth_api_key" env:"OCIS_MACHINE_AUTH_API_KEY;SEARCH_MACHINE_AUTH_API_KEY" desc:"Machine auth API key used to validate internal requests necessary for the access to resources from other services."`

	Context context.Context `yaml:"-"`
}
