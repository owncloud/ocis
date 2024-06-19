package config

import "github.com/owncloud/ocis/v2/ocis-pkg/shared"

// GRPCConfig defines the available grpc configuration.
type GRPCConfig struct {
	Addr      string                 `yaml:"addr" env:"STORE_GRPC_ADDR" desc:"The bind address of the GRPC service." introductionVersion:"pre5.0" deprecationVersion:"5.0" removalVersion:"7.0.0" deprecationInfo:"The store service is optional and will be removed."`
	Namespace string                 `yaml:"-"`
	TLS       *shared.GRPCServiceTLS `yaml:"tls"`
}
