package config

// Log defines the available log configuration.
type Log struct {
	Level  string `mapstructure:"level" env:"OCIS_LOG_LEVEL;GRAPH_EXPLORER_LOG_LEVEL" desc:"The log level."`
	Pretty bool   `mapstructure:"pretty" env:"OCIS_LOG_PRETTY;GRAPH_EXPLORER_LOG_PRETTY" desc:"Enable pretty logs."`
	Color  bool   `mapstructure:"color" env:"OCIS_LOG_COLOR;GRAPH_EXPLORER_LOG_COLOR" desc:"Enable colored logs."`
	File   string `mapstructure:"file" env:"OCIS_LOG_FILE;GRAPH_EXPLORER_LOG_FILE" desc:"The path to the log file when logging to file."`
}
