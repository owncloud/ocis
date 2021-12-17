package config

// HTTP defines the available http configuration.
type HTTP struct {
	Addr      string `ocisConfig:"addr" env:"IDP_HTTP_ADDR"`
	Root      string `ocisConfig:"root" env:"IDP_HTTP_ROOT"`
	Namespace string
	TLSCert   string `ocisConfig:"tls_cert" env:"IDP_TRANSPORT_TLS_CERT"`
	TLSKey    string `ocisConfig:"tls_key" env:"IDP_TRANSPORT_TLS_KEY"`
	TLS       bool   `ocisConfig:"tls" env:"IDP_TLS"`
}
