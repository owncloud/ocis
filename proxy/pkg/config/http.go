package config

// HTTP defines the available http configuration.
type HTTP struct {
	Addr      string `ocisConfig:"addr" env:"PROXY_HTTP_ADDR"`
	Root      string `ocisConfig:"root" env:"PROXY_HTTP_ROOT"`
	Namespace string
	TLSCert   string `ocisConfig:"tls_cert" env:"PROXY_TRANSPORT_TLS_CERT"`
	TLSKey    string `ocisConfig:"tls_key" env:"PROXY_TRANSPORT_TLS_KEY"`
	TLS       bool   `ocisConfig:"tls" env:"PROXY_TLS"`
}
