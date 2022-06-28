package config

// Log defines the available log configuration.
type Log struct {
	Level  string `mapstructure:"level" env:"OCIS_LOG_LEVEL;GRAPH_EXPLORER_LOG_LEVEL" desc:"The log level. Valid values are: \"panic\", \"fatal\", \"error\", \"warn\", \"info\", \"debug\", \"trace\"."`
	Pretty bool   `mapstructure:"pretty" env:"OCIS_LOG_PRETTY;GRAPH_EXPLORER_LOG_PRETTY" desc:"Activates pretty log output."`
	Color  bool   `mapstructure:"color" env:"OCIS_LOG_COLOR;GRAPH_EXPLORER_LOG_COLOR" desc:"Activates colorized log output."`
	File   string `mapstructure:"file" env:"OCIS_LOG_FILE;GRAPH_EXPLORER_LOG_FILE" desc:"The path to the log file. Activates logging to this file if set."`
}
