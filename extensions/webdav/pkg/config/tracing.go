package config

// Tracing defines the available tracing configuration.
type Tracing struct {
	Enabled   bool   `yaml:"enabled" env:"OCIS_TRACING_ENABLED;WEBDAV_TRACING_ENABLED" desc:"Enable tracing."`
	Type      string `yaml:"type" env:"OCIS_TRACING_TYPE;WEBDAV_TRACING_TYPE" desc:"The tracing type."`
	Endpoint  string `yaml:"endpoint" env:"OCIS_TRACING_ENDPOINT;WEBDAV_TRACING_ENDPOINT" desc:"The tracing service endpoint."`
	Collector string `yaml:"collector" env:"OCIS_TRACING_COLLECTOR;WEBDAV_TRACING_COLLECTOR" desc:"The tracing collector."`
}
