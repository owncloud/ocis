package config

// Log defines the available log configuration.
type Log struct {
	Level  string `ocisConfig:"level" env:"OCIS_LOG_LEVEL;IDP_LOG_LEVEL"`
	Pretty bool   `ocisConfig:"pretty" env:"OCIS_LOG_PRETTY;IDP_LOG_PRETTY"`
	Color  bool   `ocisConfig:"color" env:"OCIS_LOG_COLOR;IDP_LOG_COLOR"`
	File   string `ocisConfig:"file" env:"OCIS_LOG_FILE;IDP_LOG_FILE"`
}
