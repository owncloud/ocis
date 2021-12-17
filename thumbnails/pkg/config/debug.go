package config

// Debug defines the available debug configuration.
type Debug struct {
	Addr   string `ocisConfig:"addr" env:"THUMBNAILS_DEBUG_ADDR"`
	Token  string `ocisConfig:"token" env:"THUMBNAILS_DEBUG_TOKEN"`
	Pprof  bool   `ocisConfig:"pprof" env:"THUMBNAILS_DEBUG_PPROF"`
	Zpages bool   `ocisConfig:"zpages" env:"THUMBNAILS_DEBUG_ZPAGES"`
}
