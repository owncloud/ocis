package config

// Log defines the available log configuration.
type Log struct {
	Level  string `mapstructure:"level" env:"OCIS_LOG_LEVEL;WEBDAV_LOG_LEVEL" desc:"The log level."`
	Pretty bool   `mapstructure:"pretty" env:"OCIS_LOG_PRETTY;WEBDAV_LOG_PRETTY" desc:"Enable pretty log output."`
	Color  bool   `mapstructure:"color" env:"OCIS_LOG_COLOR;WEBDAV_LOG_COLOR" desc:"Enable colored log output."`
	File   string `mapstructure:"file" env:"OCIS_LOG_FILE;WEBDAV_LOG_FILE" desc:"The path to the file if the log should write to file."`
}
