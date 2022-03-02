package config

// Tracing defines the available tracing configuration.
type Tracing struct {
	Enabled   bool   `ocisConfig:"enabled" env:"OCIS_TRACING_ENABLED;IDM_TRACING_ENABLED"`
	Type      string `ocisConfig:"type" env:"OCIS_TRACING_TYPE;IDM_TRACING_TYPE"`
	Endpoint  string `ocisConfig:"endpoint" env:"OCIS_TRACING_ENDPOINT;IDM_TRACING_ENDPOINT"`
	Collector string `ocisConfig:"collector" env:"OCIS_TRACING_COLLECTOR;IDM_TRACING_COLLECTOR"`
}
