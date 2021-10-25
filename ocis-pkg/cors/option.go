package cors

import (
	"github.com/owncloud/ocis/ocis-pkg/log"
)

// Option defines a single option function.
type Option func(o *Options)

// Options defines the available options for this package.
type Options struct {
	// Logger to use for logging, must be set
	Logger log.Logger
	// AllowedOrigins represents the allowed CORS origins
	AllowedOrigins []string
	// AllowedMethods represents the allowed CORS methods
	AllowedMethods []string
	// AllowedHeaders represents the allowed CORS headers
	AllowedHeaders []string
	// AllowCredentials represents the AllowCredentials CORS option
	AllowCredentials bool
}

// newAccountOptions initializes the available default options.
func NewOptions(opts ...Option) Options {
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

// AllowedOrigins provides a function to set the AllowedOrigins option.
func AllowedOrigins(origins []string) Option {
	return func(o *Options) {
		o.AllowedOrigins = origins
	}
}

// AllowedMethods provides a function to set the AllowedMethods option.
func AllowedMethods(methods []string) Option {
	return func(o *Options) {
		o.AllowedMethods = methods
	}
}

// AllowedHeaders provides a function to set the AllowedHeaders option.
func AllowedHeaders(headers []string) Option {
	return func(o *Options) {
		o.AllowedHeaders = headers
	}
}

// AlloweCredentials provides a function to set the AllowCredentials option.
func AllowCredentials(allow bool) Option {
	return func(o *Options) {
		o.AllowCredentials = allow
	}
}
