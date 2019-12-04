package middleware

import (
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
	Config *config.Config
}

func Config(val *config.Config) Option {
	return func(o *Options) {
		o.Config = val
	}
}
