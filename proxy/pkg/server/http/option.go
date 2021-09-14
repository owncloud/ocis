package http

import (
	"context"
	"net/http"

	"github.com/justinas/alice"
	"github.com/owncloud/ocis/ocis-pkg/log"
	"github.com/owncloud/ocis/proxy/pkg/config"
	"github.com/owncloud/ocis/proxy/pkg/metrics"
	"github.com/urfave/cli/v2"
)

// Option defines a single option function.
type Option func(o *Options)

// Options defines the available options for this package.
type Options struct {
	Logger      log.Logger
	Context     context.Context
	Config      *config.Config
	Handler     http.Handler
	Metrics     *metrics.Metrics
	Flags       []cli.Flag
	Middlewares alice.Chain
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

// Handler provides a function to set the Handler option.
func Handler(h http.Handler) Option {
	return func(o *Options) {
		o.Handler = h
	}
}

// Middlewares provides a function to register middlewares
func Middlewares(val alice.Chain) Option {
	return func(o *Options) {
		o.Middlewares = val
	}
}
