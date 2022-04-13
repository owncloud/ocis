package config

// GRPC defines the available grpc configuration.
type GRPC struct {
	Addr      string `yaml:"addr" env:"THUMBNAILS_GRPC_ADDR"`
	Namespace string `yaml:"-"`
}
