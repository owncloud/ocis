package config

import (
	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
)

// HTTP defines the available http configuration.
type HTTP struct {
	Addr      string                `yaml:"addr" env:"COLLABORATION_HTTP_ADDR" desc:"The address of the HTTP service."`
	BindAddr  string                `yaml:"bindaddr" env:"COLLABORATION_HTTP_BINDADDR" desc:"The bind address of the HTTP service."`
	Namespace string                `yaml:"-"`
	Scheme    string                `yaml:"scheme" env:"COLLABORATION_HTTP_SCHEME" desc:"Either http or https"`
	TLS       shared.HTTPServiceTLS `yaml:"tls"`
}
