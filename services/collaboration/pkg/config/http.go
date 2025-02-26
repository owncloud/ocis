package config

import (
	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
)

// HTTP defines the available http configuration.
type HTTP struct {
	Addr      string                `yaml:"addr" env:"COLLABORATION_HTTP_ADDR" desc:"The bind address of the HTTP service." introductionVersion:"6.0.0"`
	Namespace string                `yaml:"-"`
	TLS       shared.HTTPServiceTLS `yaml:"tls"`
}
