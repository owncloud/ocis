package service

import (
	gatewayv1beta1 "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/collaboration/pkg/config"
	microstore "go-micro.dev/v4/store"
)

// Option defines a single option function.
type Option func(o *Options)

// Options defines the available options for this package.
type Options struct {
	Logger  log.Logger
	Config  *config.Config
	AppURLs map[string]map[string]string
	Gwc     gatewayv1beta1.GatewayAPIClient
	Store   microstore.Store
}

// newOptions initializes the available default options.
func newOptions(opts ...Option) Options {
	opt := Options{}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}

// Logger provides a function to set the Logger option.
func Logger(val log.Logger) Option {
	return func(o *Options) {
		o.Logger = val
	}
}

// Config provides a function to set the Config option.
func Config(val *config.Config) Option {
	return func(o *Options) {
		o.Config = val
	}
}

// AppURLs provides a function to set the AppURLs option.
func AppURLs(val map[string]map[string]string) Option {
	return func(o *Options) {
		o.AppURLs = val
	}
}

// GatewayAPIClient provides a function to set the GatewayAPIClient option.
func GatewayAPIClient(val gatewayv1beta1.GatewayAPIClient) Option {
	return func(o *Options) {
		o.Gwc = val
	}
}

// Store proivdes a function to set the store
func Store(val microstore.Store) Option {
	return func(o *Options) {
		o.Store = val
	}
}
