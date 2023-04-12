package oidc

import (
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
)

// Option defines a single option function.
type Option func(o *Options)

// Options defines the available options for this package.
type Options struct {
	// Logger to use for logging, must be set
	Logger log.Logger
	// The OpenID Connect Issuer URL
	OidcIssuer string
}

// newOptions initializes the available default options.
func newOptions(opts ...Option) Options {
	opt := Options{}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}

// WithLogger provides a function to set the openid connect issuer option.
func WithOidcIssuer(val string) Option {
	return func(o *Options) {
		o.OidcIssuer = val
	}
}

// WithLogger provides a function to set the logger option.
func WithLogger(val log.Logger) Option {
	return func(o *Options) {
		o.Logger = val
	}
}
