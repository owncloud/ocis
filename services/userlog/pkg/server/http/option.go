package http

import (
	"context"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	ehsvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/eventhistory/v0"
	settingssvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/settings/v0"
	"github.com/owncloud/ocis/v2/services/userlog/pkg/config"
	"github.com/owncloud/ocis/v2/services/userlog/pkg/metrics"
	"github.com/owncloud/reva/v2/pkg/events"
	"github.com/owncloud/reva/v2/pkg/rgrpc/todo/pool"
	"github.com/urfave/cli/v2"
	"go-micro.dev/v4/store"
	"go.opentelemetry.io/otel/trace"
)

// Option defines a single option function.
type Option func(o *Options)

// Options defines the available options for this package.
type Options struct {
	Logger           log.Logger
	Context          context.Context
	Config           *config.Config
	Metrics          *metrics.Metrics
	Flags            []cli.Flag
	Namespace        string
	Store            store.Store
	Stream           events.Stream
	GatewaySelector  pool.Selectable[gateway.GatewayAPIClient]
	HistoryClient    ehsvc.EventHistoryService
	ValueClient      settingssvc.ValueService
	RoleClient       settingssvc.RoleService
	RegisteredEvents []events.Unmarshaller
	TracerProvider   trace.TracerProvider
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

// Namespace provides a function to set the Namespace option.
func Namespace(val string) Option {
	return func(o *Options) {
		o.Namespace = val
	}
}

// Store provides a function to configure the store
func Store(store store.Store) Option {
	return func(o *Options) {
		o.Store = store
	}
}

// Stream provides a function to configure the stream
func Stream(stream events.Stream) Option {
	return func(o *Options) {
		o.Stream = stream
	}
}

// GatewaySelector provides a function to configure the gateway client selector
func GatewaySelector(gatewaySelector pool.Selectable[gateway.GatewayAPIClient]) Option {
	return func(o *Options) {
		o.GatewaySelector = gatewaySelector
	}
}

// History provides a function to configure the event history client
func History(h ehsvc.EventHistoryService) Option {
	return func(o *Options) {
		o.HistoryClient = h
	}
}

// RegisteredEvents provides a function to register events
func RegisteredEvents(evs []events.Unmarshaller) Option {
	return func(o *Options) {
		o.RegisteredEvents = evs
	}
}

// Value provides a function to configure the value service client
func Value(vs settingssvc.ValueService) Option {
	return func(o *Options) {
		o.ValueClient = vs
	}
}

// Roles provides a function to configure the roles service client
func Role(rs settingssvc.RoleService) Option {
	return func(o *Options) {
		o.RoleClient = rs
	}
}

// TracerProvider provides a function to set the TracerProvider option
func TracerProvider(val trace.TracerProvider) Option {
	return func(o *Options) {
		o.TracerProvider = val
	}
}
