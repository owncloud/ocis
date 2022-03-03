package config

// Log defines the available log configuration.
type Log struct {
	Level  string `ocisConfig:"level" env:"OCIS_LOG_LEVEL;ACCOUNTS_LOG_LEVEL" desc:"The log level."`
	Pretty bool   `ocisConfig:"pretty" env:"OCIS_LOG_PRETTY;ACCOUNTS_LOG_PRETTY" desc:"Activates pretty log output."`
	Color  bool   `ocisConfig:"color" env:"OCIS_LOG_COLOR;ACCOUNTS_LOG_COLOR" desc:"Activates colorized log output."`
	File   string `ocisConfig:"file" env:"OCIS_LOG_FILE;ACCOUNTS_LOG_FILE" desc:"The target log file."`
}
