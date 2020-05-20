package runtime

import (
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis-pkg/v2/log"
)

// Options is a runtime option
type Options struct {
	Services []string
	Logger   log.Logger
	Context  *cli.Context
}

// Option undocummented
type Option func(o *Options)

// newOptions initializes the available default options.
func newOptions(opts ...Option) Options {
	opt := Options{}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}

// Services option
func Services(s []string) Option {
	return func(o *Options) {
		o.Services = append(o.Services, s...)
	}
}

// Context option
func Context(c *cli.Context) Option {
	return func(o *Options) {
		o.Context = c
	}
}
