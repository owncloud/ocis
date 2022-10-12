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

	OcisPublicURL   string          `yaml:"ocis_public_url" env:"OCIS_URL;OCIS_PUBLIC_URL" desc:"URL, where oCIS is reachable for users."`
	WebdavNamespace string          `yaml:"webdav_namespace" env:"WEBDAV_WEBDAV_NAMESPACE" desc:"CS3 path layout to use when forwarding /webdav requests"`
	Reva            shared.Reva     `yaml:"reva"`
	Context         context.Context `yaml:"-"`
}
