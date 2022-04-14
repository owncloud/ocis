package config

// HTTP defines the available http configuration.
type HTTP struct {
	Addr      string `yaml:"addr" env:"IDP_HTTP_ADDR"`
	Root      string `yaml:"root" env:"IDP_HTTP_ROOT"`
	Namespace string `yaml:"-"`
	TLSCert   string `yaml:"tls_cert" env:"IDP_TRANSPORT_TLS_CERT"`
	TLSKey    string `yaml:"tls_key" env:"IDP_TRANSPORT_TLS_KEY"`
	TLS       bool   `yaml:"tls" env:"IDP_TLS"`
}
