package config

// Tracing defines the available tracing configuration.
type Tracing struct {
	Enabled   bool   `ocisConfig:"enabled" env:"OCIS_TRACING_ENABLED;NATS_TRACING_ENABLED"`
	Type      string `ocisConfig:"type" env:"OCIS_TRACING_TYPE;NATS_TRACING_TYPE"`
	Endpoint  string `ocisConfig:"endpoint" env:"OCIS_TRACING_ENDPOINT;NATS_TRACING_ENDPOINT"`
	Collector string `ocisConfig:"collector" env:"OCIS_TRACING_COLLECTOR;NATS_TRACING_COLLECTOR"`
}
