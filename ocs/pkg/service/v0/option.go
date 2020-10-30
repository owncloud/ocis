package svc

import (
	gatewayv1beta1 "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	"github.com/cs3org/reva/pkg/token"
	"net/http"

	"github.com/owncloud/ocis/ocis-pkg/log"
	"github.com/owncloud/ocis/ocs/pkg/config"
)

// Option defines a single option function.
type Option func(o *Options)

// Options defines the available options for this package.
type Options struct {
	Logger       log.Logger
	Config       *config.Config
	Middleware   []func(http.Handler) http.Handler
	TokenManager token.Manager
	RevaClient   gatewayv1beta1.GatewayAPIClient
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

// Config provides a function to set the config option.
func Config(val *config.Config) Option {
	return func(o *Options) {
		o.Config = val
	}
}

// Middleware provides a function to set the middleware option.
func Middleware(val ...func(http.Handler) http.Handler) Option {
	return func(o *Options) {
		o.Middleware = val
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
