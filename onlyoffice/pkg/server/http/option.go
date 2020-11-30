package http

import (
	"context"

	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis-pkg/v2/log"
	"github.com/owncloud/ocis/onlyoffice/pkg/config"
	"github.com/owncloud/ocis/onlyoffice/pkg/metrics"
)

// Option defines a single option function.
type Option func(o *Options)

// Options defines the available options for this package.
type Options struct {
	Name    string
	Logger  log.Logger
	Context context.Context
	Config  *config.Config
	Metrics *metrics.Metrics
	Flags   []cli.Flag
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

// Metrics provides a function to set the metrics option.
func Metrics(val *metrics.Metrics) Option {
	return func(o *Options) {
		o.Metrics = val
	}
}

// Flags provides a function to set the flags option.
func Flags(val []cli.Flag) Option {
	return func(o *Options) {
		o.Flags = append(o.Flags, val...)
	}
}

// Name provides a function to set the Name option.
func Name(val string) Option {
	return func(o *Options) {
		o.Name = val
	}
}
