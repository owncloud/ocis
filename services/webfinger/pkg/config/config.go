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

	Instances []Instance `yaml:"instances"`
	Relations []string   `yaml:"relations" env:"WEBFINGER_RELATIONS" desc:"A comma-separated list of relation URIs or registered relation types to add to webfinger responses."`
	IDP       string     `yaml:"idp" env:"OCIS_URL;WEBFINGER_OIDC_ISSUER" desc:"The identity provider href for the openid-discovery relation."`
	OcisURL   string     `yaml:"idp" env:"OCIS_URL;WEBFINGER_OCIS_URL" desc:"The oCIS instance URL for the owncloud instance relations."`

	Context context.Context `yaml:"-"`
}

// Instance to use with a matching rule and titles
type Instance struct {
	Claim  string            `yaml:"claim"`
	Regex  string            `yaml:"rule"`
	Href   string            `yaml:"href"`
	Titles map[string]string `yaml:"title"`
	Break  bool              `yaml:"break"`
}
