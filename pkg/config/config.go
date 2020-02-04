// Package config should be moved to internal
package config

// Server configures a server.
type Server struct {
	Name      string
	Namespace string
	Address   string
}

// Config merges all Account config parameters.
type Config struct {
	MountPath string
	Manager   string
	Server    Server
}

// New returns a new config.
func New() *Config {
	return &Config{}
}
