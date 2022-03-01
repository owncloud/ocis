package config

// Log defines the available log configuration.
type Log struct {
	Level  string `env:"OCIS_LOG_LEVEL;ACCOUNTS_LOG_LEVEL" desc:"The log level."`
	Pretty bool   `env:"OCIS_LOG_PRETTY;ACCOUNTS_LOG_PRETTY" desc:"Activates pretty log output."`
	Color  bool   `env:"OCIS_LOG_COLOR;ACCOUNTS_LOG_COLOR" desc:"Activates colorized log output."`
	File   string `env:"OCIS_LOG_FILE;ACCOUNTS_LOG_FILE" desc:"The target log file."`
}
