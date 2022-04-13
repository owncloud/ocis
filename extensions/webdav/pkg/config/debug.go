package config

// Debug defines the available debug configuration.
type Debug struct {
	Addr   string `yaml:"addr" env:"WEBDAV_DEBUG_ADDR"`
	Token  string `yaml:"token" env:"WEBDAV_DEBUG_TOKEN"`
	Pprof  bool   `yaml:"pprof" env:"WEBDAV_DEBUG_PPROF"`
	Zpages bool   `yaml:"zpages" env:"WEBDAV_DEBUG_ZPAGES"`
}
