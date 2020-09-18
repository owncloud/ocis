package account

import (
	"github.com/owncloud/ocis/ocis-pkg/log"
)

// Option defines a single option function.
type Option func(o *Options)

// Options defines the available options for this package.
type Options struct {
	// Logger to use for logging, must be set
	Logger log.Logger
	// JWTSecret is the jwt secret for the reva token manager
	JWTSecret string
}

// Logger provides a function to set the logger option.
func Logger(l log.Logger) Option {
	return func(o *Options) {
		o.Logger = l
	}
}

// JWTSecret provides a function to set the jwt secret option.
func JWTSecret(s string) Option {
	return func(o *Options) {
		o.JWTSecret = s
	}
}
