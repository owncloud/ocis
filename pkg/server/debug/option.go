package debug

import (
	"context"

	"github.com/owncloud/ocis-pkg/v2/log"
	"github.com/owncloud/ocis-reva/pkg/config"
)

// Option defines a single option function.
type Option func(o *Options)

// Options defines the available options for this package.
type Options struct {
	Name    string
	Addr    string
	Logger  log.Logger
	Context context.Context
	Config  *config.Config
}

// newOptions initializes the available default options.
func newOptions(opts ...Option) Options {
	opt := Options{}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}

// Name provides a function to set the name option.
func Name(val string) Option {
	return func(o *Options) {
		o.Name = val
	}
}

// Addr provides a function to set the addr option.
func Addr(val string) Option {
	return func(o *Options) {
		o.Addr = val
	}
}

// Logger provides a function to set the logger option.
func Logger(val log.Logger) Option {
	return func(o *Options) {
		o.Logger = val
	}
}

// Context provides a function to set the context option.
func Context(val context.Context) Option {
	return func(o *Options) {
		o.Context = val
	}
}

// Config provides a function to set the config option.
func Config(val *config.Config) Option {
	return func(o *Options) {
		o.Config = val
	}
}
