package config

import "github.com/owncloud/ocis/v2/ocis-pkg/shared"

// GRPCConfig defines the available grpc configuration.
type GRPCConfig struct {
	Addr      string                 `yaml:"addr" env:"SEARCH_GRPC_ADDR" desc:"The bind address of the GRPC service." introductionVersion:"pre5.0"`
	Namespace string                 `yaml:"-"`
	TLS       *shared.GRPCServiceTLS `yaml:"tls"`
}
