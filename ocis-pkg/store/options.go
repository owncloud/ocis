package store

import (
	"context"
	"time"

	"go-micro.dev/v4/store"
)

// CacheOptions are cache specific options to configure the store
type CacheOptions struct {
	// Type determines the implementation:
	// * "noop", for a noop store (it does nothing)
	// * "etcd", for etcd
	// * "ocmem", custom in-memory implementation, with fixed size and optimized prefix
	// and suffix search
	// * "memory", for a in-memory implementation, which is the default if noone matches
	Type string

	// Size configures the maximum capacity of the cache for
	// the "ocmem" implementation, in number of items that the cache can hold per table.
	// You can use 5000 to make the cache hold up to 5000 elements.
	// The parameter only affects to the "ocmem" implementation, the rest will ignore it.
	// If an invalid value is used, the default of 512 will be used instead.
	Size int

	// TTL is the time to live for documents stored in the store
	TTL time.Duration
}

type cacheOptionsContextKey struct{}

// CacheOptions defines cache options when using a store.
func WithCacheOptions(cacheOptions CacheOptions) store.Option {
	return func(o *store.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}

		o.Context = context.WithValue(o.Context, cacheOptionsContextKey{}, cacheOptions)
	}
}
