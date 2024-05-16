package config

import "github.com/owncloud/ocis/v2/ocis-pkg/shared"

// HTTP defines the available http configuration.
type HTTP struct {
	Addr                  string                `yaml:"addr" env:"THUMBNAILS_HTTP_ADDR" desc:"The bind address of the HTTP service." introductionVersion:"pre5.0"`
	TLS                   shared.HTTPServiceTLS `yaml:"tls"`
	Root                  string                `yaml:"root" env:"THUMBNAILS_HTTP_ROOT" desc:"Subdirectory that serves as the root for this HTTP service." introductionVersion:"pre5.0"`
	Namespace             string                `yaml:"-"`
	MaxConcurrentRequests int                   `yaml:"max_concurrent_requests" env:"THUMBNAILS_MAX_CONCURRENT_REQUESTS" desc:"Number of maximum concurrent thumbnail requests. Default is 0 which is unlimited." introductionVersion:"6.0"`
}
