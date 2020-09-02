package service

import (
	"github.com/owncloud/ocis-accounts/pkg/config"
	"github.com/owncloud/ocis-pkg/v2/log"
	"github.com/owncloud/ocis-pkg/v2/roles"
	settings "github.com/owncloud/ocis-settings/pkg/proto/v0"
)

// Option defines a single option function.
type Option func(o *Options)

// Options defines the available options for this package.
type Options struct {
	Logger      log.Logger
	Config      *config.Config
	RoleService settings.RoleService
	RoleManager *roles.Manager
}

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

// RoleService provides a function to set the role service option.
func RoleService(val settings.RoleService) Option {
	return func(o *Options) {
		o.RoleService = val
	}
}

// RoleManager provides a function to set the roles manager option.
func RoleManager(val *roles.Manager) Option {
	return func(o *Options) {
		o.RoleManager = val
	}
}
