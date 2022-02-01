package config

// Log defines the available log configuration.
type Log struct {
	Level  string `env:"OCIS_LOG_LEVEL;ACCOUNTS_LOG_LEVEL"`
	Pretty bool   `env:"OCIS_LOG_PRETTY;ACCOUNTS_LOG_PRETTY"`
	Color  bool   `env:"OCIS_LOG_COLOR;ACCOUNTS_LOG_COLOR"`
	File   string `env:"OCIS_LOG_FILE;ACCOUNTS_LOG_FILE"`
}
