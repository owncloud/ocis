package config

// HTTP defines the available http configuration.
type HTTP struct {
	Addr      string `yaml:"addr" env:"IDP_HTTP_ADDR" desc:"The bind address of the HTTP service."`
	Root      string `yaml:"root" env:"IDP_HTTP_ROOT" desc:"The root path of the HTTP service."`
	Namespace string `yaml:"-"`
	TLSCert   string `yaml:"tls_cert" env:"IDP_TRANSPORT_TLS_CERT"`
	TLSKey    string `yaml:"tls_key" env:"IDP_TRANSPORT_TLS_KEY"`
	TLS       bool   `yaml:"tls" env:"IDP_TLS"`
}
