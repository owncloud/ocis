package config

// Tracing defines the available tracing configuration.
type Tracing struct {
	Enabled   bool   `yaml:"enabled" env:"OCIS_TRACING_ENABLED;IDP_TRACING_ENABLED"`
	Type      string `yaml:"type" env:"OCIS_TRACING_TYPE;IDP_TRACING_TYPE" desc:"The type of tracing. Defaults to \"\", which is the same as \"jaeger\". Allowed tracing types are \"jaeger\" and \"\" as of now."`
	Endpoint  string `yaml:"endpoint" env:"OCIS_TRACING_ENDPOINT;IDP_TRACING_ENDPOINT" desc:"The endpoint of the tracing agent."`
	Collector string `yaml:"collector" env:"OCIS_TRACING_COLLECTOR;IDP_TRACING_COLLECTOR" desc:"The HTTP endpoint for sending spans directly to a collector, i.e. http://jaeger-collector:14268/api/traces. Only used if the tracing endpoint is unset."`
}
