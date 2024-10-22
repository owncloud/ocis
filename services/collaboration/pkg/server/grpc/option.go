package grpc

import (
	"context"

	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/collaboration/pkg/config"
	microstore "go-micro.dev/v4/store"
	"go.opentelemetry.io/otel/trace"
)

// Option defines a single option function.
type Option func(o *Options)

// Options defines the available options for this package.
type Options struct {
	AppURLs       map[string]map[string]string
	Name          string
	Logger        log.Logger
	Context       context.Context
	Config        *config.Config
	TraceProvider trace.TracerProvider
	Store         microstore.Store
}

// newOptions initializes the available default options.
func newOptions(opts ...Option) Options {
	opt := Options{}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}

// AppURLs provides app urls based on mimetypes.
func AppURLs(val map[string]map[string]string) Option {
	return func(o *Options) {
		o.AppURLs = val
	}
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

// TraceProvider provides a function to set the trace provider option.
func TraceProvider(val trace.TracerProvider) Option {
	return func(o *Options) {
		o.TraceProvider = val
	}
}

// Store provides a funtion to set the Store option
func Store(val microstore.Store) Option {
	return func(o *Options) {
		o.Store = val
	}
}
