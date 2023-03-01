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

	TokenManager *TokenManager `yaml:"token_manager"`

	Context context.Context `yaml:"-"`
}

// Instance to use with a matching rule and titles
type Instance struct {
	Claim  string            `yaml:"claim"`
	Regex  string            `yaml:"regex"`
	Href   string            `yaml:"href"`
	Titles map[string]string `yaml:"titles"`
	Break  bool              `yaml:"break"`
}
