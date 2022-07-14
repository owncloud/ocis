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

	GraphExplorer GraphExplorer `yaml:"graph_explorer"`

	Context context.Context `yaml:"-"`
}

// GraphExplorer defines the available graph-explorer configuration.
type GraphExplorer struct {
	ClientID     string `yaml:"client_id" env:"GRAPH_EXPLORER_CLIENT_ID" desc:"OIDC client ID the graph explorer uses. This client needs to be set up in your IDP."`
	Issuer       string `yaml:"issuer" env:"OCIS_URL;OCIS_OIDC_ISSUER;GRAPH_EXPLORER_ISSUER" desc:"URL of the OIDC issuer. It defaults to URL of the builtin IDP."`
	GraphURLBase string `yaml:"graph_url_base" env:"OCIS_URL;GRAPH_EXPLORER_GRAPH_URL_BASE" desc:"Base URL where the graph explorer is reachable for users."`
	GraphURLPath string `yaml:"graph_url_path" env:"GRAPH_EXPLORER_GRAPH_URL_PATH" desc:"URL path where the graph explorer is reachable for users."`
}
