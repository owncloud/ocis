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

	HTTP HTTP         `yaml:"http"`
	Reva *shared.Reva `yaml:"reva"`

	Instances        []Instance `yaml:"instances"`
	InstanceSelector string     `yaml:"instance_selector" env:"WEBFINGER_INSTANCE_SELECTOR" desc:"How to select which instance to use for an account. Can be 'default', 'regex' or 'claims'?"`
	InstanceLookup   string     `yaml:"instance_lookup" env:"WEBFINGER_INSTANCE_LOOKUP" desc:"How to look up to instance href and topic. Can be 'default', 'template', 'static' or 'ldap'?"`
	InstanceMatches  string     `yaml:"instance_matches" env:"WEBFINGER_INSTANCE_MATCHES" desc:"TODO"`
	LookupChain      string     `yaml:"lookup_chain" env:"WEBFINGER_LOOKUP_CHAIN" desc:"A chain of lookup steps for webfinger."`
	IDP              string     `yaml:"idp" env:"OCIS_URL;WEBFINGER_OIDC_ISSUER" desc:"The identity provider href for the openid-discovery relation."`
	OcisURL          string     `yaml:"idp" env:"OCIS_URL;WEBFINGER_OCIS_URL" desc:"The oCIS instance URL for the owncloud-account relation. The host part will be used for the 'acct' URI."`

	Rules string `yaml:"webdav_namespace" env:"WEBFINGER_" desc:"Jail requests to /dav/webdav into this CS3 namespace. Supports template layouting with CS3 User properties."`
	// TODO wie proxy?

	Context context.Context `yaml:"-"`
}

// Instance to use with a matching rule and titles
type Instance struct {
	Claim  string            `yaml:"claim"`
	Regex  string            `yaml:"rule"`
	Href   string            `yaml:"href"`
	Titles map[string]string `yaml:"title"`
}
