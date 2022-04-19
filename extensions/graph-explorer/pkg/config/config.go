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

	GraphExplorer GraphExplorer `yaml:"graph_explorer,omitempty"`

	Context context.Context `yaml:"-"`
}

// GraphExplorer defines the available graph-explorer configuration.
type GraphExplorer struct {
	ClientID     string `yaml:"client_id" env:"GRAPH_EXPLORER_CLIENT_ID"`
	Issuer       string `yaml:"issuer" env:"OCIS_URL;GRAPH_EXPLORER_ISSUER"`
	GraphURLBase string `yaml:"graph_url_base" env:"OCIS_URL;GRAPH_EXPLORER_GRAPH_URL_BASE"`
	GraphURLPath string `yaml:"graph_url_path" env:"GRAPH_EXPLORER_GRAPH_URL_PATH"`
}
