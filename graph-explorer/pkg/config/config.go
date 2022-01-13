package config

import (
	"context"

	"github.com/owncloud/ocis/ocis-pkg/shared"
)

// Config combines all available configuration parts.
type Config struct {
	*shared.Commons

	Service Service

	Tracing *Tracing `ocisConfig:"tracing"`
	Log     *Log     `ocisConfig:"log"`
	Debug   Debug    `ocisConfig:"debug"`

	HTTP HTTP `ocisConfig:"http"`

	GraphExplorer GraphExplorer `ocisConfig:"graph_explorer"`

	Context context.Context
}

// GraphExplorer defines the available graph-explorer configuration.
type GraphExplorer struct {
	ClientID     string `ocisConfig:"client_id" env:"GRAPH_EXPLORER_CLIENT_ID"`
	Issuer       string `ocisConfig:"issuer" env:"OCIS_URL;GRAPH_EXPLORER_ISSUER"`
	GraphURLBase string `ocisConfig:"graph_url_base" env:"OCIS_URL;GRAPH_EXPLORER_GRAPH_URL_BASE"`
	GraphURLPath string `ocisConfig:"graph_url_path" env:"GRAPH_EXPLORER_GRAPH_URL_PATH"`
}
