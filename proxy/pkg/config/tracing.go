package config

// Tracing defines the available tracing configuration.
type Tracing struct {
	Enabled   bool   `yaml:"enabled" env:"OCIS_TRACING_ENABLED;PROXY_TRACING_ENABLED"`
	Type      string `yaml:"type" env:"OCIS_TRACING_TYPE;PROXY_TRACING_TYPE"`
	Endpoint  string `yaml:"endpoint" env:"OCIS_TRACING_ENDPOINT;PROXY_TRACING_ENDPOINT"`
	Collector string `yaml:"collector" env:"OCIS_TRACING_COLLECTOR;PROXY_TRACING_COLLECTOR"`
}
