package config

// GRPC defines the available grpc configuration.
type GRPC struct {
	Addr      string `yaml:"addr" env:"COLLABORATION_GRPC_ADDR" desc:"The bind address of the GRPC service." introductionVersion:"%%NEXT%%"`
	Namespace string `yaml:"-"`
}
