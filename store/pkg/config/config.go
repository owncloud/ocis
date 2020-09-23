package config

// Log defines the available logging configuration.
type Log struct {
	Level  string
	Pretty bool
	Color  bool
}

// Debug defines the available debug configuration.
type Debug struct {
	Addr   string
	Token  string
	Pprof  bool
	Zpages bool
}

// GRPC defines the available grpc configuration.
type GRPC struct {
	Addr      string
	Namespace string
	Root      string
}

// Tracing defines the available tracing configuration.
type Tracing struct {
	Enabled   bool
	Type      string
	Endpoint  string
	Collector string
	Service   string
}

// Config combines all available configuration parts.
type Config struct {
	File     string
	Log      Log
	Debug    Debug
	GRPC     GRPC
	Tracing  Tracing
	Datapath string
	Name     string
}

// New initializes a new configuration with or without defaults.
func New() *Config {
	return &Config{}
}
