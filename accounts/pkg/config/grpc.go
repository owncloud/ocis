package config

// GRPC defines the available grpc configuration.
type GRPC struct {
	Addr      string `env:"ACCOUNTS_GRPC_ADDR"`
	Namespace string
}
