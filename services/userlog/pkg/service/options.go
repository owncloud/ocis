package service

import (
	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	"github.com/go-chi/chi/v5"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	ehsvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/eventhistory/v0"
	settingssvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/settings/v0"
	"github.com/owncloud/ocis/v2/services/userlog/pkg/config"
	"github.com/owncloud/reva/v2/pkg/events"
	"github.com/owncloud/reva/v2/pkg/rgrpc/todo/pool"
	"go-micro.dev/v4/store"
	"go.opentelemetry.io/otel/trace"
)

// Option for the userlog service
type Option func(*Options)

// Options for the userlog service
type Options struct {
	Logger           log.Logger
	Stream           events.Stream
	Mux              *chi.Mux
	Store            store.Store
	Config           *config.Config
	HistoryClient    ehsvc.EventHistoryService
	GatewaySelector  pool.Selectable[gateway.GatewayAPIClient]
	ValueClient      settingssvc.ValueService
	RoleClient       settingssvc.RoleService
	RegisteredEvents []events.Unmarshaller
	TraceProvider    trace.TracerProvider
}

// Logger configures a logger for the userlog service
func Logger(log log.Logger) Option {
	return func(o *Options) {
		o.Logger = log
	}
}

// Stream configures an event stream for the userlog service
func Stream(s events.Stream) Option {
	return func(o *Options) {
		o.Stream = s
	}
}

// Mux defines the muxer for the userlog service
func Mux(m *chi.Mux) Option {
	return func(o *Options) {
		o.Mux = m
	}
}

// Store defines the store for the userlog service
func Store(s store.Store) Option {
	return func(o *Options) {
		o.Store = s
	}
}

// Config adds the config for the userlog service
func Config(c *config.Config) Option {
	return func(o *Options) {
		o.Config = c
	}
}

// HistoryClient adds a grpc client for the eventhistory service
func HistoryClient(hc ehsvc.EventHistoryService) Option {
	return func(o *Options) {
		o.HistoryClient = hc
	}
}

// GatewaySelector adds a grpc client selector for the gateway service
func GatewaySelector(gatewaySelector pool.Selectable[gateway.GatewayAPIClient]) Option {
	return func(o *Options) {
		o.GatewaySelector = gatewaySelector
	}
}

// RegisteredEvents registers the events the service should listen to
func RegisteredEvents(e []events.Unmarshaller) Option {
	return func(o *Options) {
		o.RegisteredEvents = e
	}
}

// ValueClient adds a grpc client for the value service
func ValueClient(vs settingssvc.ValueService) Option {
	return func(o *Options) {
		o.ValueClient = vs
	}
}

// RoleClient adds a grpc client for the role service
func RoleClient(rs settingssvc.RoleService) Option {
	return func(o *Options) {
		o.RoleClient = rs
	}
}

// TraceProvider adds a tracer provider for the userlog service
func TraceProvider(tp trace.TracerProvider) Option {
	return func(o *Options) {
		o.TraceProvider = tp
	}
}
