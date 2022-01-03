package config

import (
	"context"

	"github.com/owncloud/ocis/ocis-pkg/shared"
)

// Config combines all available configuration parts.
type Config struct {
	*shared.Commons

	Service Service

	Tracing Tracing `ocisConfig:"tracing"`
	Log     *Log    `ocisConfig:"log"`
	Debug   Debug   `ocisConfig:"debug"`

	HTTP HTTP `ocisConfig:"http"`

	OcisPublicURL   string `ocisConfig:"ocis_public_url" env:"OCIS_URL;OCIS_PUBLIC_URL"`
	WebdavNamespace string `ocisConfig:"webdav_namespace" env:"STORAGE_WEBDAV_NAMESPACE"`

	Context context.Context
}
