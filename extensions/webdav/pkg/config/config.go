package config

import (
	"context"

	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
)

// Config combines all available configuration parts.
type Config struct {
	Commons *shared.Commons `yaml:"-"` // don't use this directly as configuration for a service

	Service Service `yaml:"-"`

	Tracing *Tracing `yaml:"tracing"`
	Log     *Log     `yaml:"log"`
	Debug   Debug    `yaml:"debug"`

	HTTP HTTP `yaml:"http"`

	OcisPublicURL   string `yaml:"ocis_public_url" env:"OCIS_URL;OCIS_PUBLIC_URL"`
	WebdavNamespace string `yaml:"webdav_namespace" env:"WEBDAV_WEBDAV_NAMESPACE" desc:"CS3 path layout to use when forwarding /webdav requests"` //TODO: prevent this cross config
	RevaGateway     string `yaml:"reva_gateway" env:"REVA_GATEWAY" desc:"The CS3 gateway endpoint"`

	Context context.Context `yaml:"-"`
}
