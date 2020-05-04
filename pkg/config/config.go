// Package config should be moved to internal
package config

// Server configures a server.
type Server struct {
	Name      string
	Namespace string
	Address   string
}

// Log defines the available logging configuration.
type Log struct {
	Level  string
	Pretty bool
	Color  bool
}

// Config merges all Account config parameters.
type Config struct {
	MountPath string
	Manager   string
	Server    Server
	Log       Log
}

// New returns a new config.
func New() *Config {
	return &Config{}
}
