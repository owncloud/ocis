package config

// Log defines the available log configuration.
type Log struct {
	Level  string `mapstructure:"level" env:"OCIS_LOG_LEVEL;NATS_LOG_LEVEL"`
	Pretty bool   `mapstructure:"pretty" env:"OCIS_LOG_PRETTY;NATS_LOG_PRETTY"`
	Color  bool   `mapstructure:"color" env:"OCIS_LOG_COLOR;NATS_LOG_COLOR"`
	File   string `mapstructure:"file" env:"OCIS_LOG_FILE;NATS_LOG_FILE"`
}
