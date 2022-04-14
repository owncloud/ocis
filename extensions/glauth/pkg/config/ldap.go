package config

// Ldap defines the available LDAP configuration.
type Ldap struct {
	Enabled   bool   `yaml:"enabled" env:"GLAUTH_LDAP_ENABLED"`
	Addr      string `yaml:"addr" env:"GLAUTH_LDAP_ADDR"`
	Namespace string `yaml:"-"`
}
