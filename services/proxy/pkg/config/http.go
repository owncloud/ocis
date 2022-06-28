package config

// HTTP defines the available http configuration.
type HTTP struct {
	Addr      string `yaml:"addr" env:"PROXY_HTTP_ADDR" desc:"The bind address of the HTTP service."`
	Root      string `yaml:"root" env:"PROXY_HTTP_ROOT" desc:"The root path of the HTTP service."`
	Namespace string `yaml:"-"`
	TLSCert   string `yaml:"tls_cert" env:"PROXY_TRANSPORT_TLS_CERT"`
	TLSKey    string `yaml:"tls_key" env:"PROXY_TRANSPORT_TLS_KEY"`
	TLS       bool   `yaml:"tls" env:"PROXY_TLS"`
}
