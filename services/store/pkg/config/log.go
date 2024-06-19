package config

// Log defines the available log configuration.
type Log struct {
	Level  string `mapstructure:"level" env:"OCIS_LOG_LEVEL;STORE_LOG_LEVEL" desc:"The log level. Valid values are: 'panic', 'fatal', 'error', 'warn', 'info', 'debug', 'trace'." introductionVersion:"pre5.0" deprecationVersion:"5.0" removalVersion:"7.0.0" deprecationInfo:"The store service is optional and will be removed."`
	Pretty bool   `mapstructure:"pretty" env:"OCIS_LOG_PRETTY;STORE_LOG_PRETTY" desc:"Activates pretty log output." introductionVersion:"pre5.0" deprecationVersion:"5.0" removalVersion:"7.0.0" deprecationInfo:"The store service is optional and will be removed."`
	Color  bool   `mapstructure:"color" env:"OCIS_LOG_COLOR;STORE_LOG_COLOR" desc:"Activates colorized log output." introductionVersion:"pre5.0" deprecationVersion:"5.0" removalVersion:"7.0.0" deprecationInfo:"The store service is optional and will be removed."`
	File   string `mapstructure:"file" env:"OCIS_LOG_FILE;STORE_LOG_FILE" desc:"The path to the log file. Activates logging to this file if set." introductionVersion:"pre5.0" deprecationVersion:"5.0" removalVersion:"7.0.0" deprecationInfo:"The store service is optional and will be removed."`
}
