package config

// CORS defines the available cors configuration.
type CORS struct {
	AllowedOrigins   []string `ocisConfig:"allowed_origins"`
	AllowedMethods   []string `ocisConfig:"allowed_methods"`
	AllowedHeaders   []string `ocisConfig:"allowed_headers"`
	AllowCredentials bool     `ocisConfig:"allow_credentials"`
}

// HTTP defines the available http configuration.
type HTTP struct {
	Addr      string `ocisConfig:"addr" env:"WEBDAV_HTTP_ADDR"`
	Namespace string
	Root      string `ocisConfig:"root" env:"WEBDAV_HTTP_ROOT"`
	CORS      CORS   `ocisConfig:"cors"`
}
