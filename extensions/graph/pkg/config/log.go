package config

// Log defines the available log configuration.
type Log struct {
	Level  string `mapstructure:"level" env:"OCIS_LOG_LEVEL;GRAPH_LOG_LEVEL"`
	Pretty bool   `mapstructure:"pretty" env:"OCIS_LOG_PRETTY;GRAPH_LOG_PRETTY"`
	Color  bool   `mapstructure:"color" env:"OCIS_LOG_COLOR;GRAPH_LOG_COLOR"`
	File   string `mapstructure:"file" env:"OCIS_LOG_FILE;GRAPH_LOG_FILE"`
}
