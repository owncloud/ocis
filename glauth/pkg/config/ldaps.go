package config

// Ldaps defined the available LDAPS configuration.
type Ldaps struct {
	Enabled   bool   `ocisConfig:"enabled" env:"GLAUTH_LDAPS_ENABLED"`
	Addr      string `ocisConfig:"addr" env:"GLAUTH_LDAPS_ADDR"`
	Namespace string
	Cert      string `ocisConfig:"cert" env:"GLAUTH_LDAPS_CERT"`
	Key       string `ocisConfig:"key" env:"GLAUTH_LDAPS_KEY"`
}
