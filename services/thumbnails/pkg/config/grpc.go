package config

import "github.com/owncloud/ocis/v2/ocis-pkg/shared"

// GRPCConfig defines the available grpc configuration.
type GRPCConfig struct {
	Addr                  string                 `yaml:"addr" env:"THUMBNAILS_GRPC_ADDR" desc:"The bind address of the GRPC service." introductionVersion:"pre5.0"`
	Namespace             string                 `yaml:"-"`
	TLS                   *shared.GRPCServiceTLS `yaml:"tls"`
	MaxConcurrentRequests int                    `yaml:"max_concurrent_requests" env:"THUMBNAILS_MAX_CONCURRENT_REQUESTS" desc:"Number of maximum concurrent thumbnail requests. Default is 0 which is unlimited." introductionVersion:"6.0.0"`
}
