package config

// Log defines the available log configuration.
type Log struct {
	Level  string `mapstructure:"level" env:"OCIS_LOG_LEVEL;CLIENTLOG_USERLOG_LOG_LEVEL" desc:"The log level. Valid values are: 'panic', 'fatal', 'error', 'warn', 'info', 'debug', 'trace'." introductionVersion:"5.0"`
	Pretty bool   `mapstructure:"pretty" env:"OCIS_LOG_PRETTY;CLIENTLOG_USERLOG_LOG_PRETTY" desc:"Activates pretty log output." introductionVersion:"5.0"`
	Color  bool   `mapstructure:"color" env:"OCIS_LOG_COLOR;CLIENTLOG_USERLOG_LOG_COLOR" desc:"Activates colorized log output." introductionVersion:"5.0"`
	File   string `mapstructure:"file" env:"OCIS_LOG_FILE;CLIENTLOG_USERLOG_LOG_FILE" desc:"The path to the log file. Activates logging to this file if set." introductionVersion:"5.0"`
}
