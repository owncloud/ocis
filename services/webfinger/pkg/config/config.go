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

	InstanceSelector string `yaml:"instance_selector" env:"WEBFINGER_INSTANCE_SELECTOR" desc:"How to select which instance to use for an account. Can be 'default', 'regex' or 'claims'?"`
	InstanceLookup   string `yaml:"instance_lookup" env:"WEBFINGER_INSTANCE_LOOKUP" desc:"How to look up to instance href and topic. Can be 'default', 'template', 'static' or 'ldap'?"`

	Rules string `yaml:"webdav_namespace" env:"WEBFINGER_" desc:"Jail requests to /dav/webdav into this CS3 namespace. Supports template layouting with CS3 User properties."`
	// TODO wie proxy?

	Context context.Context `yaml:"-"`
}
