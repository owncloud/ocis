package svc

import (
	"net/http"

	"github.com/owncloud/ocis/ocis-pkg/log"
	"github.com/owncloud/ocis/ocis-pkg/roles"
	"github.com/owncloud/ocis/ocs/pkg/config"
	settings "github.com/owncloud/ocis/settings/pkg/proto/v0"
)

// Option defines a single option function.
type Option func(o *Options)

// Options defines the available options for this package.
type Options struct {
	Logger      log.Logger
	Config      *config.Config
	Middleware  []func(http.Handler) http.Handler
	RoleService settings.RoleService
	RoleManager *roles.Manager
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

// RoleService provides a function to set the RoleService option.
func RoleService(val settings.RoleService) Option {
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
