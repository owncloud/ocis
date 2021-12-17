package config

// Debug defines the available debug configuration.
type Debug struct {
	Addr   string `ocisConfig:"addr" env:"IDP_DEBUG_ADDR"`
	Token  string `ocisConfig:"token" env:"IDP_DEBUG_TOKEN"`
	Pprof  bool   `ocisConfig:"pprof" env:"IDP_DEBUG_PPROF"`
	Zpages bool   `ocisConfig:"zpages" env:"IDP_DEBUG_ZPAGES"`
}
