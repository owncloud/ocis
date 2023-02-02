package config

// HTTP defines the available http configuration.
type HTTP struct {
	Addr      string `yaml:"addr" env:"PROXY_HTTP_ADDR" desc:"The bind address of the HTTP service."`
	Root      string `yaml:"root" env:"PROXY_HTTP_ROOT" desc:"Subdirectory that serves as the root for this HTTP service."`
	Namespace string `yaml:"-"`
	TLSCert   string `yaml:"tls_cert" env:"PROXY_TRANSPORT_TLS_CERT" desc:"Path/File name of the TLS server certificate (in PEM format) for the external http services. If not definied, the root directory derives from $OCIS_BASE_DATA_PATH:/proxy."`
	TLSKey    string `yaml:"tls_key" env:"PROXY_TRANSPORT_TLS_KEY" desc:"Path/File name for the TLS certificate key (in PEM format) for the server certificate to use for the external http services. If not definied, the root directory derives from $OCIS_BASE_DATA_PATH:/proxy."`
	TLS       bool   `yaml:"tls" env:"PROXY_TLS" desc:"Enable/Disable HTTPS for the external HTTP services."`
}
