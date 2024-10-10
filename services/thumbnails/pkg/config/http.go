package config

import "github.com/owncloud/ocis/v2/ocis-pkg/shared"

// HTTP defines the available http configuration.
type HTTP struct {
	Addr      string                `yaml:"addr" env:"THUMBNAILS_HTTP_ADDR" desc:"The bind address of the HTTP service." introductionVersion:"pre5.0"`
	TLS       shared.HTTPServiceTLS `yaml:"tls"`
	Root      string                `yaml:"root" env:"THUMBNAILS_HTTP_ROOT" desc:"Subdirectory that serves as the root for this HTTP service." introductionVersion:"pre5.0"`
	Namespace string                `yaml:"-"`
}
