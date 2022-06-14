package config

// Tracing defines the available tracing configuration.
type Tracing struct {
	Enabled   bool   `ocisConfig:"enabled" env:"OCIS_TRACING_ENABLED;SEARCH_TRACING_ENABLED" desc:"Activates tracing."`
	Type      string `ocisConfig:"type" env:"OCIS_TRACING_TYPE;SEARCH_TRACING_TYPE" desc:"The type of tracing. Defaults to \"\", which is the same as \"jaeger\". Allowed tracing types are \"jaeger\" and \"\" as of now."`
	Endpoint  string `ocisConfig:"endpoint" env:"OCIS_TRACING_ENDPOINT;SEARCH_TRACING_ENDPOINT" desc:"The endpoint of the tracing agent."`
	Collector string `ocisConfig:"collector" env:"OCIS_TRACING_COLLECTOR;SEARCH_TRACING_COLLECTOR" desc:"The HTTP endpoint for sending spans directly to a collector, i.e. http://jaeger-collector:14268/api/traces. Only used if the tracing endpoint is unset."`
}
