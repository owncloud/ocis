package grpc

import (
	"context"

	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"go-micro.dev/v4/server"
	"go.opentelemetry.io/otel/trace"
)

// Option defines a single option function.
type Option func(o *Options)

// Options defines the available options for this package.
type Options struct {
	Logger          log.Logger
	Namespace       string
	Name            string
	Version         string
	Address         string
	TLSEnabled      bool
	TLSCert         string
	TLSKey          string
	Context         context.Context
	TraceProvider   trace.TracerProvider
	HandlerWrappers []server.HandlerWrapper
}

// newOptions initializes the available default options.
func newOptions(opts ...Option) Options {
	opt := Options{
		Namespace: "go.micro.api",
	}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}

// Logger provides a function to set the logger option.
func Logger(l log.Logger) Option {
	return func(o *Options) {
		o.Logger = l
	}
}

// Namespace provides a function to set the namespace option.
func Namespace(n string) Option {
	return func(o *Options) {
		o.Namespace = n
	}
}

// Name provides a function to set the name option.
func Name(n string) Option {
	return func(o *Options) {
		o.Name = n
	}
}

// Version provides a function to set the version option.
func Version(v string) Option {
	return func(o *Options) {
		o.Version = v
	}
}

// Address provides a function to set the address option.
func Address(a string) Option {
	return func(o *Options) {
		o.Address = a
	}
}

// TLSEnabled provides a function to enable/disable TLS
func TLSEnabled(v bool) Option {
	return func(o *Options) {
		o.TLSEnabled = v
	}
}

// TLSCert provides a function to set the TLS server certificate and key
func TLSCert(c string, k string) Option {
	return func(o *Options) {
		o.TLSCert = c
		o.TLSKey = k
	}
}

// Context provides a function to set the context option.
func Context(ctx context.Context) Option {
	return func(o *Options) {
		o.Context = ctx
	}
}

// TraceProvider provides a function to set the trace provider option.
func TraceProvider(tp trace.TracerProvider) Option {
	return func(o *Options) {
		o.TraceProvider = tp
	}
}

func HandlerWrappers(w ...server.HandlerWrapper) Option {
	return func(o *Options) {
		o.HandlerWrappers = w
	}
}
