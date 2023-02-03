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

	Rules string `yaml:"webdav_namespace" env:"WEBFINGER_" desc:"Jail requests to /dav/webdav into this CS3 namespace. Supports template layouting with CS3 User properties."`
	// TODO wie proxy?

	Context context.Context `yaml:"-"`
}
