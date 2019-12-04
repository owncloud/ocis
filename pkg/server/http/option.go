package http

import (
	"context"

	"github.com/owncloud/ocis-graph/pkg/config"
)

func newOptions(opts ...Option) Options {
	opt := Options{}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}

type Option func(o *Options)

type Options struct {
	Context context.Context
	Config  *config.Config
}

func Context(val context.Context) Option {
	return func(o *Options) {
		o.Context = val
	}
}

func Config(val *config.Config) Option {
	return func(o *Options) {
		o.Config = val
	}
}
