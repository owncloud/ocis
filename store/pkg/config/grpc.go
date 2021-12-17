package config

// GRPC defines the available grpc configuration.
type GRPC struct {
	Addr      string `ocisConfig:"addr" env:"STORE_GRPC_ADDR"`
	Namespace string
}
