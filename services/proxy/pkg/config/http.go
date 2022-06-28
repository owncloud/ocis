package config

// HTTP defines the available http configuration.
type HTTP struct {
	Addr      string `yaml:"addr" env:"PROXY_HTTP_ADDR"`
	Root      string `yaml:"root" env:"PROXY_HTTP_ROOT"`
	Namespace string `yaml:"-"`
	TLSCert   string `yaml:"tls_cert" env:"PROXY_TRANSPORT_TLS_CERT"`
	TLSKey    string `yaml:"tls_key" env:"PROXY_TRANSPORT_TLS_KEY"`
	TLS       bool   `yaml:"tls" env:"PROXY_TLS"`
}
