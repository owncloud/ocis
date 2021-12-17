package config

// Tracing defines the available tracing configuration.
type Tracing struct {
	Enabled   bool   `ocisConfig:"enabled" env:"OCIS_TRACING_ENABLED;ACCOUNTS_TRACING_ENABLED"`
	Type      string `ocisConfig:"type" env:"OCIS_TRACING_TYPE;ACCOUNTS_TRACING_TYPE"`
	Endpoint  string `ocisConfig:"endpoint" env:"OCIS_TRACING_ENDPOINT;ACCOUNTS_TRACING_ENDPOINT"`
	Collector string `ocisConfig:"collector" env:"OCIS_TRACING_COLLECTOR;ACCOUNTS_TRACING_COLLECTOR"`
	Service   string `ocisConfig:"service" env:"ACCOUNTS_TRACING_SERVICE"` //TODO: should this be an ID? or the same as Service.Name?
}
