package grpc

import (
	"context"

	"github.com/owncloud/ocis/ocis-pkg/log"
	"github.com/owncloud/ocis/thumbnails/pkg/config"
	"github.com/owncloud/ocis/thumbnails/pkg/metrics"
	"github.com/urfave/cli/v2"
)

// Option defines a single option function.
type Option func(o *Options)

// Options defines the available options for this package.
type Options struct {
	Name      string
	Address   string
	Logger    log.Logger
	Context   context.Context
	Config    *config.Config
	Metrics   *metrics.Metrics
	Namespace string
	Flags     []cli.Flag
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

// Name provides a name for the service.
func Name(val string) Option {
	return func(o *Options) {
		o.Name = val
	}
}

// Address provides an address for the service.
func Address(val string) Option {
	return func(o *Options) {
		o.Address = val
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

// Namespace provides a function to set the namespace option.
func Namespace(val string) Option {
	return func(o *Options) {
		o.Namespace = val
	}
}

// Flags provides a function to set the flags option.
func Flags(flags []cli.Flag) Option {
	return func(o *Options) {
		o.Flags = append(o.Flags, flags...)
	}
}
