package config

// Log defines the available log configuration.
type Log struct {
	Level  string `mapstructure:"level" env:"OCIS_LOG_LEVEL;PROXY_LOG_LEVEL"`
	Pretty bool   `mapstructure:"pretty" env:"OCIS_LOG_PRETTY;PROXY_LOG_PRETTY"`
	Color  bool   `mapstructure:"color" env:"OCIS_LOG_COLOR;PROXY_LOG_COLOR"`
	File   string `mapstructure:"file" env:"OCIS_LOG_FILE;PROXY_LOG_FILE"`
}
