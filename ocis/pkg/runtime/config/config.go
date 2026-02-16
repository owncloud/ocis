package config

// Config determines behavior across the tool.
type Config struct {
	// Hostname where the runtime is running. When using PMAN in cli mode, it determines where the host runtime is.
	// Default is localhost.
	Hostname string

	// Port configures the port where a runtime is available. It defaults to 10666.
	Port string

	// KeepAlive configures if restart attempts are made if the process supervised terminates. Default is false.
	KeepAlive bool
}

var (
	defaultHostname = "localhost"
	defaultPort     = "10666"
)

// NewConfig returns a new config with a set of defaults.
func NewConfig() *Config {
	return &Config{
		Hostname:  defaultHostname,
		Port:      defaultPort,
		KeepAlive: false,
	}
}
