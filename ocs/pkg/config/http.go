package config

// HTTP defines the available http configuration.
type HTTP struct {
	Addr      string `ocisConfig:"addr" env:"OCS_HTTP_ADDR"`
	Root      string `ocisConfig:"root" env:"OCS_HTTP_ROOT"`
	Namespace string
	CORS      CORS `ocisConfig:"cors"`
}

// CORS defines the available cors configuration.
type CORS struct {
	AllowedOrigins   []string `ocisConfig:"allowed_origins"`
	AllowedMethods   []string `ocisConfig:"allowed_methods"`
	AllowedHeaders   []string `ocisConfig:"allowed_headers"`
	AllowCredentials bool     `ocisConfig:"allowed_credentials"`
}
