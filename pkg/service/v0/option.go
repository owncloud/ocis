package service

import (
	"github.com/owncloud/ocis-pkg/v2/log"
	"github.com/owncloud/ocis-store/pkg/config"
)

// Option defines a single option function.
type Option func(o *Options)

// Options defines the available options for this package.
type Options struct {
	Logger log.Logger
	Config *config.Config

	Database, Table string
	Nodes           []string
}

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

func Database(val *config.Config) Option {
	return func(o *Options) {
		o.Config = val
	}
}

func Table(val *config.Config) Option {
	return func(o *Options) {
		o.Config = val
	}
}

func Nodes(val *config.Config) Option {
	return func(o *Options) {
		o.Config = val
	}
}

func Config(val *config.Config) Option {
	return func(o *Options) {
		o.Config = val
	}
}
