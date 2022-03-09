package config

// Debug defines the available debug configuration.
type Debug struct {
	Addr   string `ocisConfig:"addr" env:"AUDIT_DEBUG_ADDR"`
	Token  string `ocisConfig:"token" env:"AUDIT_DEBUG_TOKEN"`
	Pprof  bool   `ocisConfig:"pprof" env:"AUDIT_DEBUG_PPROF"`
	Zpages bool   `ocisConfig:"zpages" env:"AUDIT_DEBUG_ZPAGES"`
}
