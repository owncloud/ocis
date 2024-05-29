package config

import (
	"context"

	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
	"github.com/owncloud/ocis/v2/ocis-pkg/tracing"
)

// Config combines all available configuration parts.
type Config struct {
	Commons *shared.Commons `yaml:"-"` // don't use this directly as configuration for a service

	Service Service `yaml:"-"`
	App     App     `yaml:"app"`

	TokenManager *TokenManager `yaml:"token_manager"`

	GRPC GRPC `yaml:"grpc"`
	HTTP HTTP `yaml:"http"`

	Wopi   Wopi   `yaml:"wopi"`
	CS3Api CS3Api `yaml:"cs3api"`

	Tracing *Tracing `yaml:"tracing"`
	Log     *Log     `yaml:"log"`
	Debug   Debug    `yaml:"debug"`

	Context context.Context `yaml:"-"`
}

// Tracing defines the available tracing configuration. Not used at the moment
type Tracing struct {
	Enabled   bool   `yaml:"enabled" env:"OCIS_TRACING_ENABLED;COLLABORATION_TRACING_ENABLED" desc:"Activates tracing." introductionVersion:"5.1"`
	Type      string `yaml:"type" env:"OCIS_TRACING_TYPE;COLLABORATION_TRACING_TYPE" desc:"The type of tracing. Defaults to '', which is the same as 'jaeger'. Allowed tracing types are 'jaeger' and '' as of now." introductionVersion:"5.1"`
	Endpoint  string `yaml:"endpoint" env:"OCIS_TRACING_ENDPOINT;COLLABORATION_TRACING_ENDPOINT" desc:"The endpoint of the tracing agent." introductionVersion:"5.1"`
	Collector string `yaml:"collector" env:"OCIS_TRACING_COLLECTOR;COLLABORATION_TRACING_COLLECTOR" desc:"The HTTP endpoint for sending spans directly to a collector, i.e. http://jaeger-collector:14268/api/traces. Only used if the tracing endpoint is unset." introductionVersion:"5.1"`
}

// Convert Tracing to the tracing package's Config struct.
func (t Tracing) Convert() tracing.Config {
	return tracing.Config{
		Enabled:   t.Enabled,
		Type:      t.Type,
		Endpoint:  t.Endpoint,
		Collector: t.Collector,
	}
}
