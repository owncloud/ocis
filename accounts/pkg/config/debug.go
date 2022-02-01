package config

// Debug defines the available debug configuration.
type Debug struct {
	Addr   string `env:"ACCOUNTS_DEBUG_ADDR"`
	Token  string `env:"ACCOUNTS_DEBUG_TOKEN"`
	Pprof  bool   `env:"ACCOUNTS_DEBUG_PPROF"`
	Zpages bool   `env:"ACCOUNTS_DEBUG_ZPAGES"`
}
