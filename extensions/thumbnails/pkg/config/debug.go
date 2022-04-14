package config

// Debug defines the available debug configuration.
type Debug struct {
	Addr   string `yaml:"addr" env:"THUMBNAILS_DEBUG_ADDR"`
	Token  string `yaml:"token" env:"THUMBNAILS_DEBUG_TOKEN"`
	Pprof  bool   `yaml:"pprof" env:"THUMBNAILS_DEBUG_PPROF"`
	Zpages bool   `yaml:"zpages" env:"THUMBNAILS_DEBUG_ZPAGES"`
}
