package config

// Log defines the available log configuration.
type Log struct {
	Level  string `yaml:"level" env:"OCIS_LOG_LEVEL;COLLABORATION_LOG_LEVEL" desc:"The log level. Valid values are: 'panic', 'fatal', 'error', 'warn', 'info', 'debug', 'trace'."`
	Pretty bool   `yaml:"pretty" env:"OCIS_LOG_PRETTY;COLLABORATION_LOG_PRETTY" desc:"Activates pretty log output."`
	Color  bool   `yaml:"color" env:"OCIS_LOG_COLOR;COLLABORATION_LOG_COLOR" desc:"Activates colorized log output."`
	File   string `yaml:"file" env:"OCIS_LOG_FILE;COLLABORATION_LOG_FILE" desc:"The path to the log file. Activates logging to this file if set."`
}
