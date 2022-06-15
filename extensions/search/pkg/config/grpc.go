package config

// GRPC defines the available grpc configuration.
type GRPC struct {
	Addr      string `ocisConfig:"addr" env:"SEARCH_GRPC_ADDR" desc:"The bind address of the GRPC service."`
	Namespace string `ocisConfig:"-" yaml:"-"`
}
