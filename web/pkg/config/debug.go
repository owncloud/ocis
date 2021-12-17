package config

// Debug defines the available debug configuration.
type Debug struct {
	Addr   string `ocisConfig:"addr" env:"WEB_DEBUG_ADDR"`
	Token  string `ocisConfig:"token" env:"WEB_DEBUG_TOKEN"`
	Pprof  bool   `ocisConfig:"pprof" env:"WEB_DEBUG_PPROF"`
	Zpages bool   `ocisConfig:"zpages" env:"WEB_DEBUG_ZPAGES"`
}
