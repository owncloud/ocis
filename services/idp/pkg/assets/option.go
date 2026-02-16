package assets

import (
	"net/http"

	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/ocis-pkg/x/io/fsx"
	"github.com/owncloud/ocis/v2/services/idp"
	"github.com/owncloud/ocis/v2/services/idp/pkg/config"
)

// New returns a new http filesystem to serve assets.
func New(opts ...Option) http.FileSystem {
	options := newOptions(opts...)

	var assetFS fsx.FS = fsx.NewBasePathFs(fsx.FromIOFS(idp.Assets), "assets")

	// only use a fsx.NewFallbackFS and fsx.OsFs if a path is set, use the embedded fs only otherwise
	if options.Config.Asset.Path != "" {
		assetFS = fsx.NewFallbackFS(fsx.NewBasePathFs(fsx.NewOsFs(), options.Config.Asset.Path), assetFS)
	}

	return http.FS(assetFS.IOFS())
}

// Option defines a single option function.
type Option func(o *Options)

// Options define the available options for this package.
type Options struct {
	Logger log.Logger
	Config *config.Config
}

// newOptions initializes the available default options.
func newOptions(opts ...Option) Options {
	opt := Options{}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}

// Logger provides a function to set the logger option.
func Logger(val log.Logger) Option {
	return func(o *Options) {
		o.Logger = val
	}
}

// Config provides a function to set the config option.
func Config(val *config.Config) Option {
	return func(o *Options) {
		o.Config = val
	}
}
