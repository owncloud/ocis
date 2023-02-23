package http

import (
	"context"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	ehsvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/eventhistory/v0"
	"github.com/owncloud/ocis/v2/services/userlog/pkg/config"
	"github.com/owncloud/ocis/v2/services/userlog/pkg/metrics"
	"github.com/urfave/cli/v2"
	"go-micro.dev/v4/store"
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
	Consumer         events.Consumer
	GatewayClient    gateway.GatewayAPIClient
	HistoryClient    ehsvc.EventHistoryService
	RegisteredEvents []events.Unmarshaller
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

// Consumer provides a function to configure the consumer
func Consumer(consumer events.Consumer) Option {
	return func(o *Options) {
		o.Consumer = consumer
	}
}

// Gateway provides a function to configure the gateway client
func Gateway(gw gateway.GatewayAPIClient) Option {
	return func(o *Options) {
		o.GatewayClient = gw
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
