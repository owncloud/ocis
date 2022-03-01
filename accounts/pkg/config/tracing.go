package config

// Tracing defines the available tracing configuration.
type Tracing struct {
	Enabled   bool   `env:"OCIS_TRACING_ENABLED;ACCOUNTS_TRACING_ENABLED" desc:"Activates tracing."`
	Type      string `env:"OCIS_TRACING_TYPE;ACCOUNTS_TRACING_TYPE"`
	Endpoint  string `env:"OCIS_TRACING_ENDPOINT;ACCOUNTS_TRACING_ENDPOINT" desc:"The endpoint to the tracing collector."`
	Collector string `env:"OCIS_TRACING_COLLECTOR;ACCOUNTS_TRACING_COLLECTOR" `
}
