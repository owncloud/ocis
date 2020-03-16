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

// HTTP defines the available http configuration.
type HTTP struct {
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

// Asset defines the available asset configuration.
type Asset struct {
	Path string
}

// Policy enables us to use multiple directors.
type Policy struct {
	Name   string  `mapstructure:"name"`
	Routes []Route `mapstructure:"routes"`
}

// Route define forwarding routes
type Route struct {
	Endpoint    string `mapstructure:"endpoint"`
	Backend     string `mapstructure:"backend"`
	ApacheVHost bool   `mapstructure:"apache-vhost"`
}

// Config combines all available configuration parts.
type Config struct {
	File     string
	Log      Log
	Debug    Debug
	HTTP     HTTP
	Tracing  Tracing
	Asset    Asset
	Policies []Policy `mapstructure:"policies"`
}

// New initializes a new configuration
func New() *Config {
	return &Config{}
}
