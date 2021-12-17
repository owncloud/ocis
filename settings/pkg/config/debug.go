package config

// Debug defines the available debug configuration.
type Debug struct {
	Addr   string `ocisConfig:"addr" env:"SETTINGS_DEBUG_ADDR"`
	Token  string `ocisConfig:"token" env:"SETTINGS_DEBUG_TOKEN"`
	Pprof  bool   `ocisConfig:"pprof" env:"SETTINGS_DEBUG_PPROF"`
	Zpages bool   `ocisConfig:"zpages" env:"SETTINGS_DEBUG_ZPAGES"`
}
