package config

import "github.com/owncloud/ocis/v2/ocis-pkg/tracing"

// Tracing defines the available tracing configuration.
type Tracing struct {
	Enabled   bool   `yaml:"enabled" env:"OCIS_TRACING_ENABLED;STORE_TRACING_ENABLED" desc:"Activates tracing." introductionVersion:"pre5.0" deprecationVersion:"5.0" removalVersion:"7.0.0" deprecationInfo:"The store service is optional and will be removed."`
	Type      string `yaml:"type" env:"OCIS_TRACING_TYPE;STORE_TRACING_TYPE" desc:"The type of tracing. Defaults to '', which is the same as 'jaeger'. Allowed tracing types are 'jaeger' and '' as of now." introductionVersion:"pre5.0" deprecationVersion:"5.0" removalVersion:"7.0.0" deprecationInfo:"The store service is optional and will be removed."`
	Endpoint  string `yaml:"endpoint" env:"OCIS_TRACING_ENDPOINT;STORE_TRACING_ENDPOINT" desc:"The endpoint of the tracing agent." introductionVersion:"pre5.0" deprecationVersion:"5.0" removalVersion:"7.0.0" deprecationInfo:"The store service is optional and will be removed."`
	Collector string `yaml:"collector" env:"OCIS_TRACING_COLLECTOR;STORE_TRACING_COLLECTOR" desc:"The HTTP endpoint for sending spans directly to a collector, i.e. http://jaeger-collector:14268/api/traces. Only used if the tracing endpoint is unset." introductionVersion:"pre5.0" deprecationVersion:"5.0" removalVersion:"7.0.0" deprecationInfo:"The store service is optional and will be removed."`
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
