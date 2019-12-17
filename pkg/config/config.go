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

type HTTP struct {
	Network string
	Addr    string
	Root    string // TODO do we need the http root path
}
type GRPC struct {
	Network string
	Addr    string
}

// Reva defines the available reva configuration.
type Reva struct {
	// MaxCPUs can be a number or a percentage
	MaxCPUs  string
	LogLevel string
	// Network can be tcp, udp or unix
	HTTP      HTTP
	GRPC      GRPC
	JWTSecret string
}

// AuthProvider defines the available authprovider configuration.
type AuthProvider struct {
	Provider string
	Insecure bool
}

// Tracing defines the available tracing configuration.
type Tracing struct {
	Enabled   bool
	Type      string
	Endpoint  string
	Collector string
	Service   string
}

// Asset defines the available asset configuration.
type Asset struct {
	Path string
}

// Config combines all available configuration parts.
type Config struct {
	File         string
	Log          Log
	Debug        Debug
	Reva         Reva
	AuthProvider AuthProvider
	Tracing      Tracing
	Asset        Asset
}

// New initializes a new configuration with or without defaults.
func New() *Config {
	return &Config{}
}
