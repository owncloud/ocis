// Package config should be moved to internal
package config

// Config captures ocis-accounts configuration parameters
type Config struct {
	MountPath string
	Manager   string
}

// New returns a new config
func New() *Config {
	return &Config{}
}
