package config

// Tracing defines the available tracing configuration.
type Tracing struct {
	Enabled   bool   `yaml:"enabled" env:"OCIS_TRACING_ENABLED;GRAPH_TRACING_ENABLED"`
	Type      string `yaml:"type" env:"OCIS_TRACING_TYPE;GRAPH_TRACING_TYPE"`
	Endpoint  string `yaml:"endpoint" env:"OCIS_TRACING_ENDPOINT;GRAPH_TRACING_ENDPOINT"`
	Collector string `yaml:"collector" env:"OCIS_TRACING_COLLECTOR;GRAPH_TRACING_COLLECTOR"`
}
