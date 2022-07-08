package svc

import (
	"net/http"

	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/ocis-pkg/roles"
	settingssvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/settings/v0"
	"github.com/owncloud/ocis/v2/services/graph/pkg/config"
	"github.com/owncloud/ocis/v2/services/graph/pkg/identity"
)

// Option defines a single option function.
type Option func(o *Options)

// Options defines the available options for this package.
type Options struct {
	Logger          log.Logger
	Config          *config.Config
	Middleware      []func(http.Handler) http.Handler
	GatewayClient   GatewayClient
	IdentityBackend identity.Backend
	RoleService     settingssvc.RoleService
	RoleManager     *roles.Manager
	EventsPublisher events.Publisher
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

// WithGatewayClient provides a function to set the gateway client option.
func WithGatewayClient(val GatewayClient) Option {
	return func(o *Options) {
		o.GatewayClient = val
	}
}

// WithIdentityBackend provides a function to set the IdentityBackend option.
func WithIdentityBackend(val identity.Backend) Option {
	return func(o *Options) {
		o.IdentityBackend = val
	}
}

// RoleService provides a function to set the RoleService option.
func RoleService(val settingssvc.RoleService) Option {
	return func(o *Options) {
		o.RoleService = val
	}
}

// RoleManager provides a function to set the RoleManager option.
func RoleManager(val *roles.Manager) Option {
	return func(o *Options) {
		o.RoleManager = val
	}
}

// EventsPublisher provides a function to set the EventsPublisher option.
func EventsPublisher(val events.Publisher) Option {
	return func(o *Options) {
		o.EventsPublisher = val
	}
}
