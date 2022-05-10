package config

// Tracing defines the available tracing configuration.
type Tracing struct {
	Enabled   bool   `yaml:"enabled" env:"OCIS_TRACING_ENABLED;GRAPH_EXPLORER_TRACING_ENABLED" desc:"Enable tracing."`
	Type      string `yaml:"type" env:"OCIS_TRACING_TYPE;GRAPH_EXPLORER_TRACING_TYPE" desc:"The sampler type: remote, const, probabilistic, ratelimiting (default remote). See also https://www.jaegertracing.io/docs/latest/sampling/."`
	Endpoint  string `yaml:"endpoint" env:"OCIS_TRACING_ENDPOINT;GRAPH_EXPLORER_TRACING_ENDPOINT" desc:"The endpoint of the tracing service."`
	Collector string `yaml:"collector" env:"OCIS_TRACING_COLLECTOR;GRAPH_EXPLORER_TRACING_COLLECTOR" desc:"The HTTP endpoint for sending spans directly to a collector, i.e. http://jaeger-collector:14268/api/traces. If specified, the tracing endpoint is ignored."`
}
