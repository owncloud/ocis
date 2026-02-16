package http

import (
	"context"

	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/invitations/pkg/config"
	svc "github.com/owncloud/ocis/v2/services/invitations/pkg/service/v0"
	"github.com/urfave/cli/v2"
	"go.opentelemetry.io/otel/trace"
)

// Option defines a single option function.
type Option func(o *Options)

// Options defines the available options for this package.
type Options struct {
	Name          string
	Namespace     string
	Logger        log.Logger
	Context       context.Context
	Config        *config.Config
	Flags         []cli.Flag
	Service       svc.Service
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

// Name provides a name for the service.
func Name(val string) Option {
	return func(o *Options) {
		o.Name = val
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

// Flags provides a function to set the flags option.
func Flags(val []cli.Flag) Option {
	return func(o *Options) {
		o.Flags = append(o.Flags, val...)
	}
}

// Namespace provides a function to set the namespace option.
func Namespace(val string) Option {
	return func(o *Options) {
		o.Namespace = val
	}
}

// Service provides a function to set the service option.
func Service(val svc.Service) Option {
	return func(o *Options) {
		o.Service = val
	}
}

// TraceProvider provides a function to set the trace provider
func TraceProvider(tp trace.TracerProvider) Option {
	return func(o *Options) {
		o.TraceProvider = tp
	}
}
