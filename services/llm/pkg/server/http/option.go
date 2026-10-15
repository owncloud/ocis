// services/llm/pkg/server/http/option.go
package http

import (
	"context"

	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"

	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/llm/pkg/config"
	"github.com/owncloud/ocis/v2/services/llm/pkg/ratelimit"
	"github.com/urfave/cli/v2"
)

// Option defines a single option function.
type Option func(o *Options)

// Options defines the available options for this package.
type Options struct {
	Logger        log.Logger
	Context       context.Context
	Config        *config.Config
	Flags         []cli.Flag
	Limiter       *ratelimit.Limiter
	TraceProvider trace.TracerProvider
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
	return func(o *Options) { o.Logger = val }
}

// Context provides a function to set the context option.
func Context(val context.Context) Option {
	return func(o *Options) { o.Context = val }
}

// Config provides a function to set the config option.
func Config(val *config.Config) Option {
	return func(o *Options) { o.Config = val }
}

// Limiter provides a function to set the rate limiter option.
func Limiter(val *ratelimit.Limiter) Option {
	return func(o *Options) { o.Limiter = val }
}

// TraceProvider provides a function to configure the trace provider.
func TraceProvider(traceProvider trace.TracerProvider) Option {
	return func(o *Options) {
		if traceProvider != nil {
			o.TraceProvider = traceProvider
		} else {
			o.TraceProvider = noop.NewTracerProvider()
		}
	}
}
