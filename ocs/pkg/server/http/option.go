package http

import (
	"context"
	gatewayv1beta1 "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	"github.com/cs3org/reva/pkg/token"

	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis/ocis-pkg/log"
	"github.com/owncloud/ocis/ocs/pkg/config"
	"github.com/owncloud/ocis/ocs/pkg/metrics"
)

// Option defines a single option function.
type Option func(o *Options)

// Options defines the available options for this package.
type Options struct {
	Namespace string
	Logger    log.Logger
	Context   context.Context
	Config    *config.Config
	Metrics   *metrics.Metrics
	Flags     []cli.Flag
	TokenManager token.Manager
	RevaClient gatewayv1beta1.GatewayAPIClient
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

// TokenManager provides a function to set the TokenManager option.
func TokenManager(tm token.Manager) Option {
	return func(o *Options) {
		o.TokenManager = tm
	}
}

// RevaClient provides a function to set the RevaClient option.
func RevaClient(c gatewayv1beta1.GatewayAPIClient) Option {
	return func(o *Options) {
		o.RevaClient = c
	}
}
