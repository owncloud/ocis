package roles

import (
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	ocisstore "github.com/owncloud/ocis/v2/ocis-pkg/store"
	settingssvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/settings/v0"
)

// Options are all the possible options.
type Options struct {
	storeOptions ocisstore.OcisStoreOptions
	logger       log.Logger
	roleService  settingssvc.RoleService
}

// Option mutates option
type Option func(*Options)

// Logger sets a preconfigured logger
func Logger(logger log.Logger) Option {
	return func(o *Options) {
		o.logger = logger
	}
}

// RoleService provides endpoints for fetching roles.
func RoleService(rs settingssvc.RoleService) Option {
	return func(o *Options) {
		o.roleService = rs
	}
}

// StoreOptions
// FIXME: nolint
// nolint: revive
func StoreOptions(storeOpts ocisstore.OcisStoreOptions) Option {
	return func(o *Options) {
		o.storeOptions = storeOpts
	}
}

func newOptions(opts ...Option) Options {
	o := Options{}

	for _, v := range opts {
		v(&o)
	}

	return o
}
