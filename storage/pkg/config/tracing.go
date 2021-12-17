package config

// Tracing defines the available tracing configuration.
type Tracing struct {
	Enabled   bool   `ocisConfig:"enabled" env:"OCIS_TRACING_ENABLED;STORAGE_TRACING_ENABLED"`
	Type      string `ocisConfig:"type" env:"OCIS_TRACING_TYPE;STORAGE_TRACING_TYPE"`
	Endpoint  string `ocisConfig:"endpoint" env:"OCIS_TRACING_ENDPOINT;STORAGE_TRACING_ENDPOINT"`
	Collector string `ocisConfig:"collector" env:"OCIS_TRACING_COLLECTOR;STORAGE_TRACING_COLLECTOR"`
	Service   string `ocisConfig:"service" env:"STORAGE_TRACING_SERVICE"` //TODO: should this be an ID? or the same as Service.Name?
}
