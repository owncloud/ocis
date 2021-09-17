package glauth

import (
	"context"

	"github.com/glauth/glauth/pkg/config"
	accounts "github.com/owncloud/ocis/accounts/pkg/proto/v0"
	"github.com/owncloud/ocis/ocis-pkg/log"
)

// Option defines a single option function.
type Option func(o *Options)

// Options defines the available options for this package.
type Options struct {
	Logger          log.Logger
	Context         context.Context
	LDAP            *config.LDAP
	LDAPS           *config.LDAPS
	Backend         *config.Config
	Fallback        *config.Config
	BaseDN          string
	NameFormat      string
	GroupFormat     string
	RoleBundleUUID  string
	AccountsService accounts.AccountsService
	GroupsService   accounts.GroupsService
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

// LDAP provides a function to set the LDAP option.
func LDAP(val *config.LDAP) Option {
	return func(o *Options) {
		o.LDAP = val
	}
}

// LDAPS provides a function to set the LDAPS option.
func LDAPS(val *config.LDAPS) Option {
	return func(o *Options) {
		o.LDAPS = val
	}
}

// Backend provides a function to set the backend option.
func Backend(val *config.Config) Option {
	return func(o *Options) {
		o.Backend = val
	}
}

// Fallback provides a string to set the fallback option.
func Fallback(val *config.Config) Option {
	return func(o *Options) {
		o.Fallback = val
	}
}

// BaseDN provides a string to set the BaseDN option.
func BaseDN(val string) Option {
	return func(o *Options) {
		o.BaseDN = val
	}
}

// NameFormat provides a string to set the NameFormat option.
func NameFormat(val string) Option {
	return func(o *Options) {
		o.NameFormat = val
	}
}

// GroupFormat provides a string to set the GroupFormat option.
func GroupFormat(val string) Option {
	return func(o *Options) {
		o.GroupFormat = val
	}
}

// AccountsService provides an AccountsService client to set the AccountsService option.
func AccountsService(val accounts.AccountsService) Option {
	return func(o *Options) {
		o.AccountsService = val
	}
}

// GroupsService provides an GroupsService client to set the GroupsService option.
func GroupsService(val accounts.GroupsService) Option {
	return func(o *Options) {
		o.GroupsService = val
	}
}

// RoleBundleUUID provides a role bundle UUID to make internal grpc requests.
func RoleBundleUUID(val string) Option {
	return func(o *Options) {
		o.RoleBundleUUID = val
	}
}
