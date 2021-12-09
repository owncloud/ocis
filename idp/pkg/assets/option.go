package assets

import (
	"net/http"

	"github.com/owncloud/ocis/idp"
	"github.com/owncloud/ocis/idp/pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/assetsfs"
	"github.com/owncloud/ocis/ocis-pkg/log"
)

// New returns a new http filesystem to serve assets.
func New(opts ...Option) http.FileSystem {
	options := newOptions(opts...)
	return assetsfs.New(idp.Assets, options.Config.Asset.Path, options.Logger)
}

// Option defines a single option function.
type Option func(o *Options)

// Options defines the available options for this package.
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
