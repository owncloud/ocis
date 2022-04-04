package config

import (
	"context"

	"github.com/owncloud/ocis/ocis-pkg/shared"
)

// Config combines all available configuration parts.
type Config struct {
	*shared.Commons `yaml:"-"`

	Service Service `yaml:"-"`

	Tracing *Tracing `yaml:"tracing"`
	Log     *Log     `yaml:"log"`
	Debug   Debug    `yaml:"debug"`

	HTTP HTTP `yaml:"http"`

	OcisPublicURL   string `yaml:"ocis_public_url" env:"OCIS_URL;OCIS_PUBLIC_URL"`
	WebdavNamespace string `yaml:"webdav_namespace" env:"STORAGE_WEBDAV_NAMESPACE"` //TODO: prevent this cross config
	RevaGateway     string `yaml:"reva_gateway" env:"REVA_GATEWAY"`

	Context context.Context `yaml:"-"`
}
