package config

import "context"

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

// HTTP defines the available http configuration.
type HTTP struct {
	Addr      string
	Namespace string
	Root      string
	CacheTTL  int
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
	File    string
	Log     Log
	Debug   Debug
	HTTP    HTTP
	Tracing Tracing
	Asset   Asset

	Context    context.Context
	Supervised bool
}

// Asset defines the available asset configuration.
type Asset struct {
	Path string
}

// New initializes a new configuration with or without defaults.
func New() *Config {
	return &Config{}
}
