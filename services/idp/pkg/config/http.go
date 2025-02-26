package config

// HTTP defines the available http configuration.
type HTTP struct {
	Addr      string `yaml:"addr" env:"IDP_HTTP_ADDR" desc:"The bind address of the HTTP service." introductionVersion:"pre5.0"`
	Root      string `yaml:"root" env:"IDP_HTTP_ROOT" desc:"Subdirectory that serves as the root for this HTTP service." introductionVersion:"pre5.0"`
	Namespace string `yaml:"-"`
	TLSCert   string `yaml:"tls_cert" env:"IDP_TRANSPORT_TLS_CERT" desc:"Path/File name of the TLS server certificate (in PEM format) for the IDP service. If not defined, the root directory derives from $OCIS_BASE_DATA_PATH/idp." introductionVersion:"pre5.0"`
	TLSKey    string `yaml:"tls_key" env:"IDP_TRANSPORT_TLS_KEY" desc:"Path/File name for the TLS certificate key (in PEM format) for the server certificate to use for the IDP service. If not defined, the root directory derives from $OCIS_BASE_DATA_PATH/idp." introductionVersion:"pre5.0"`
	TLS       bool   `yaml:"tls" env:"IDP_TLS" desc:"Disable or Enable HTTPS for the communication between the Proxy service and the IDP service. If set to 'true', the key and cert files need to be configured and present." introductionVersion:"pre5.0"`
}
