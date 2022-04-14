package config

// Debug defines the available debug configuration.
type Debug struct {
	Addr   string `yaml:"addr" env:"AUDIT_DEBUG_ADDR"`
	Token  string `yaml:"token" env:"AUDIT_DEBUG_TOKEN"`
	Pprof  bool   `yaml:"pprof" env:"AUDIT_DEBUG_PPROF"`
	Zpages bool   `yaml:"zpages" env:"AUDIT_DEBUG_ZPAGES"`
}
