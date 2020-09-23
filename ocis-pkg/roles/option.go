package roles

import (
	"time"

	"github.com/owncloud/ocis/ocis-pkg/log"
	settings "github.com/owncloud/ocis/settings/pkg/proto/v0"
)

// Options are all the possible options.
type Options struct {
	size        int
	ttl         time.Duration
	logger      log.Logger
	roleService settings.RoleService
}

// Option mutates option
type Option func(*Options)

// CacheSize configures the size of the cache in items.
func CacheSize(s int) Option {
	return func(o *Options) {
		o.size = s
	}
}

// CacheTTL rebuilds the cache after the configured duration.
func CacheTTL(ttl time.Duration) Option {
	return func(o *Options) {
		o.ttl = ttl
	}
}

// Logger sets a preconfigured logger
func Logger(logger log.Logger) Option {
	return func(o *Options) {
		o.logger = logger
	}
}

// RoleService provides endpoints for fetching roles.
func RoleService(rs settings.RoleService) Option {
	return func(o *Options) {
		o.roleService = rs
	}
}

func newOptions(opts ...Option) Options {
	o := Options{}

	for _, v := range opts {
		v(&o)
	}

	return o
}
