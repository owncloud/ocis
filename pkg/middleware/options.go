package middleware

import (
	"github.com/owncloud/ocis-ocs/pkg/config"
	"github.com/owncloud/ocis-pkg/v2/log"
)

// Option defines a single option function.
type Option func(o *Options)

// Options defines the available options for this package.
type Options struct {
	// Logger to use for logging, must be set
	Logger log.Logger
	// TokenManagerConfig for communicating with the reva token manager
	TokenManagerConfig config.TokenManager
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
func Logger(l log.Logger) Option {
	return func(o *Options) {
		o.Logger = l
	}
}

// TokenManagerConfig provides a function to set the token manger config option.
func TokenManagerConfig(cfg config.TokenManager) Option {
	return func(o *Options) {
		o.TokenManagerConfig = cfg
	}
}
