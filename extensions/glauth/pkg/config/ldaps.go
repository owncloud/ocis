package config

// Ldaps defined the available LDAPS configuration.
type Ldaps struct {
	Enabled   bool   `yaml:"enabled" env:"GLAUTH_LDAPS_ENABLED"`
	Addr      string `yaml:"addr" env:"GLAUTH_LDAPS_ADDR"`
	Namespace string `yaml:"-"`
	Cert      string `yaml:"cert" env:"GLAUTH_LDAPS_CERT"`
	Key       string `yaml:"key" env:"GLAUTH_LDAPS_KEY"`
}
