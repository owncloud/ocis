package config

// Debug defines the available debug configuration.
type Debug struct {
	Addr   string `ocisConfig:"addr" env:"WEBDAV_DEBUG_ADDR"`
	Token  string `ocisConfig:"token" env:"WEBDAV_DEBUG_TOKEN"`
	Pprof  bool   `ocisConfig:"pprof" env:"WEBDAV_DEBUG_PPROF"`
	Zpages bool   `ocisConfig:"zpages" env:"WEBDAV_DEBUG_ZPAGES"`
}
