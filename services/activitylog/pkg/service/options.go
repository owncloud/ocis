package service

import (
	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/activitylog/pkg/config"
	microstore "go-micro.dev/v4/store"
	"go.opentelemetry.io/otel/trace"
)

// Option for the activitylog service
type Option func(*Options)

// Options for the activitylog service
type Options struct {
	Logger           log.Logger
	Config           *config.Config
	TraceProvider    trace.TracerProvider
	Stream           events.Stream
	RegisteredEvents []events.Unmarshaller
	Store            microstore.Store
}

// Logger configures a logger for the activitylog service
func Logger(log log.Logger) Option {
	return func(o *Options) {
		o.Logger = log
	}
}

// Config adds the config for the activitylog service
func Config(c *config.Config) Option {
	return func(o *Options) {
		o.Config = c
	}
}

// TraceProvider adds a tracer provider for the activitylog service
func TraceProvider(tp trace.TracerProvider) Option {
	return func(o *Options) {
		o.TraceProvider = tp
	}
}

// Stream configures an event stream for the clientlog service
func Stream(s events.Stream) Option {
	return func(o *Options) {
		o.Stream = s
	}
}

// RegisteredEvents registers the events the service should listen to
func RegisteredEvents(e []events.Unmarshaller) Option {
	return func(o *Options) {
		o.RegisteredEvents = e
	}
}

// Store configures the store to use
func Store(store microstore.Store) Option {
	return func(o *Options) {
		o.Store = store
	}
}
