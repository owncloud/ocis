package config

// Debug defines the available debug configuration.
type Debug struct {
	Addr   string `ocisConfig:"addr" env:"GLAUTH_DEBUG_ADDR"`
	Token  string `ocisConfig:"token" env:"GLAUTH_DEBUG_TOKEN"`
	Pprof  bool   `ocisConfig:"pprof" env:"GLAUTH_DEBUG_PPROF"`
	Zpages bool   `ocisConfig:"zpages" env:"GLAUTH_DEBUG_ZPAGES"`
}
