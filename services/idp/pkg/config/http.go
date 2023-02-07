package config

// HTTP defines the available http configuration.
type HTTP struct {
	Addr      string `yaml:"addr" env:"IDP_HTTP_ADDR" desc:"The bind address of the HTTP service."`
	Root      string `yaml:"root" env:"IDP_HTTP_ROOT" desc:"Subdirectory that serves as the root for this HTTP service."`
	Namespace string `yaml:"-"`
	TLSCert   string `yaml:"tls_cert" env:"IDP_TRANSPORT_TLS_CERT" desc:"Path/File name of the TLS server certificate (in PEM format) for the IDP service. If not definied, the root directory derives from $OCIS_BASE_DATA_PATH:/idp."`
	TLSKey    string `yaml:"tls_key" env:"IDP_TRANSPORT_TLS_KEY" desc:"Path/File name for the TLS certificate key (in PEM format) for the server certificate to use for the IDP service. If not definied, the root directory derives from $OCIS_BASE_DATA_PATH:/idp."`
	TLS       bool   `yaml:"tls" env:"IDP_TLS" desc:"Enable/Disable HTTPS for the IDP service."`
}
