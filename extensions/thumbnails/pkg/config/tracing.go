package config

// Tracing defines the available tracing configuration.
type Tracing struct {
	Enabled   bool   `yaml:"enabled" env:"OCIS_TRACING_ENABLED;THUMBNAILS_TRACING_ENABLED" desc:"Enable tracing."`
	Type      string `yaml:"type" env:"OCIS_TRACING_TYPE;THUMBNAILS_TRACING_TYPE" desc:"The tracing type."`
	Endpoint  string `yaml:"endpoint" env:"OCIS_TRACING_ENDPOINT;THUMBNAILS_TRACING_ENDPOINT" desc:"The endpoint of the tracing service."`
	Collector string `yaml:"collector" env:"OCIS_TRACING_COLLECTOR;THUMBNAILS_TRACING_COLLECTOR" desc:"The tracing collector."`
}
