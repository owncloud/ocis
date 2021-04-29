package assets

import (
	graph_explorer "github.com/owncloud/ocis/graph-explorer"
	"github.com/owncloud/ocis/graph-explorer/pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/log"
	"net/http"
)

// assets gets initialized by New and provides the handler.
type assets struct {
	logger log.Logger
	config *config.Config
}

// Open just implements the HTTP filesystem interface.
func (a assets) Open(original string) (http.File, error) {
	return graph_explorer.Assets.Open(original)
}

// New returns a new http filesystem to serve assets.
func New(opts ...Option) http.FileSystem {
	options := newOptions(opts...)

	return assets{
		logger: options.Logger,
		config: options.Config,
	}
}
