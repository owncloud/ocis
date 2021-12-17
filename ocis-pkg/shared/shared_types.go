package shared

// EnvBinding represents a direct binding from an env variable to a go kind. Along with gookit/config, its primal goal
// is to unpack environment variables into a Go value. We do so with reflection, and this data structure is just a step
// in between.
type EnvBinding struct {
	EnvVars     []string    // name of the environment var.
	Destination interface{} // pointer to the original config value to modify.
}

// Log defines the available logging configuration.
type Log struct {
	Level  string `ocisConfig:"level" env:"OCIS_LOG_LEVEL"`
	Pretty bool   `ocisConfig:"pretty" env:"OCIS_LOG_PRETTY"`
	Color  bool   `ocisConfig:"color" env:"OCIS_LOG_COLOR"`
	File   string `ocisConfig:"file" env:"OCIS_LOG_FILE"`
}

// Tracing defines the available tracing configuration.
type Tracing struct {
	Enabled   bool   `ocisConfig:"enabled" env:"OCIS_TRACING_ENABLED"`
	Type      string `ocisConfig:"type" env:"OCIS_TRACING_TYPE"`
	Endpoint  string `ocisConfig:"endpoint" env:"OCIS_TRACING_ENDPOINT"`
	Collector string `ocisConfig:"collector" env:"OCIS_TRACING_COLLECTOR"`
}

// Commons holds configuration that are common to all extensions. Each extension can then decide whether
// to overwrite its values.
type Commons struct {
	Log     *Log     `ocisConfig:"log"`
	Tracing *Tracing `ocisConfig:"tracing"`
	OcisURL string   `ocisConfig:"ocis_url" env:"OCIS_URL"`
}
