package config

// Ldap defines the available LDAP configuration.
type Ldap struct {
	Enabled   bool   `ocisConfig:"enabled" env:"GLAUTH_LDAP_ENABLED"`
	Addr      string `ocisConfig:"addr" env:"GLAUTH_LDAP_ADDR"`
	Namespace string
}
