package natsjsregistry

import (
	"context"
	"time"

	"go-micro.dev/v4/registry"
	"go-micro.dev/v4/store"
)

type storeOptionsKey struct{}
type expiryKey struct{}

// StoreOptions sets the options for the underlying store
func StoreOptions(opts []store.Option) registry.Option {
	return func(o *registry.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, storeOptionsKey{}, opts)
	}
}

// ServiceExpiry allows setting an expiry time for service registrations
func ServiceExpiry(t time.Duration) registry.Option {
	return func(o *registry.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, expiryKey{}, t)
	}
}
