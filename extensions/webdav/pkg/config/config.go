package config

import (
	"context"

	"github.com/owncloud/ocis/ocis-pkg/shared"
)

// Config combines all available configuration parts.
type Config struct {
	*shared.Commons `yaml:"-"`

	Service Service `yaml:"-"`

	Tracing *Tracing `yaml:"tracing,omitempty"`
	Log     *Log     `yaml:"log,omitempty"`
	Debug   Debug    `yaml:"debug,omitempty"`

	HTTP HTTP `yaml:"http,omitempty"`

	OcisPublicURL   string `yaml:"ocis_public_url,omitempty" env:"OCIS_URL;OCIS_PUBLIC_URL"`
	WebdavNamespace string `yaml:"webdav_namespace,omitempty" env:"STORAGE_WEBDAV_NAMESPACE"` //TODO: prevent this cross config
	RevaGateway     string `yaml:"reva_gateway,omitempty" env:"REVA_GATEWAY"`

	Context context.Context `yaml:"-,omitempty"`
}
