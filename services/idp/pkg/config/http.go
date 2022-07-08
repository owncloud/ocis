package config

// HTTP defines the available http configuration.
type HTTP struct {
	Addr      string `yaml:"addr" env:"IDP_HTTP_ADDR" desc:"The bind address of the HTTP service."`
	Root      string `yaml:"root" env:"IDP_HTTP_ROOT" desc:"Subdirectory that serves as the root for this HTTP service."`
	Namespace string `yaml:"-"`
	TLSCert   string `yaml:"tls_cert" env:"IDP_TRANSPORT_TLS_CERT" desc:"File name of the TLS server certificate for the HTTPS server."`
	TLSKey    string `yaml:"tls_key" env:"IDP_TRANSPORT_TLS_KEY" desc:"File name of the TLS server certificate key for the HTTPS server."`
	TLS       bool   `yaml:"tls" env:"IDP_TLS" desc:"Use HTTPS server instead of HTTP server."`
}
