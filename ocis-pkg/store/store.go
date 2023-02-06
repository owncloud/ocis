package store

import (
	"context"
	"time"

	natsjs "github.com/go-micro/plugins/v4/store/nats-js"
	"github.com/go-micro/plugins/v4/store/redis"
	"github.com/nats-io/nats.go"
	"github.com/owncloud/ocis/v2/ocis-pkg/store/etcd"
	"github.com/owncloud/ocis/v2/ocis-pkg/store/memory"
	"go-micro.dev/v4/store"
)

var ocMemStore *store.Store

// Options are the options to configure the store
type Options struct {
	// Type determines the implementation:
	// * "noop", for a noop store (it does nothing)
	// * "etcd", for etcd
	// * "ocmem", custom in-memory implementation, with fixed size and optimized prefix
	// and suffix search
	// * "memory", for a in-memory implementation, which is the default if noone matches
	Type string

	// Address is a list of nodes that the store will use.
	Addresses []string

	// Size configures the maximum capacity of the cache for
	// the "ocmem" implementation, in number of items that the cache can hold per table.
	// You can use 5000 to make the cache hold up to 5000 elements.
	// The parameter only affects to the "ocmem" implementation, the rest will ignore it.
	// If an invalid value is used, the default of 512 will be used instead.
	Size int

	// Database the store should use (optional)
	Database string

	// Table the store should use (optional)
	Table string

	// TTL is the time to life for documents stored in the store
	TTL time.Duration
}

// Create returns a configured key-value store
//
// Each microservice (or whatever piece is using the store) should use the
// options available in the interface's operations to choose the right database
// and table to prevent collisions with other microservices.
// Recommended approach is to use "services" or "ocis-pkg" for the database,
// and "services/<service-name>/" or "ocis-pkg/<pkg>/" for the package name.
func Create(opts ...Option) store.Store {
	options := &Options{}
	for _, o := range opts {
		o(options)
	}

	storeopts := storeOptions(options)

	switch options.Type {
	default:
		// TODO: better to error in default case?
		fallthrough
	case "mem":
		return store.NewMemoryStore(storeopts...)
	case "noop":
		return store.NewNoopStore(storeopts...)
	case "etcd":
		return etcd.NewEtcdStore(storeopts...)
	case "redis":
		// FIXME redis plugin does not support redis cluster, sentinel or ring -> needs upstream patch or our implementation
		return redis.NewStore(storeopts...)
	case "ocmem":
		if ocMemStore == nil {
			var memStore store.Store

			sizeNum := options.Size
			if sizeNum <= 0 {
				memStore = memory.NewMultiMemStore()
			} else {
				memStore = memory.NewMultiMemStore(
					store.WithContext(
						memory.NewContext(
							context.Background(),
							map[string]interface{}{
								"maxCap": sizeNum,
							},
						)),
				)
			}
			ocMemStore = &memStore
		}
		return *ocMemStore
	case "nats-js":
		// TODO nats needs a DefaultTTL option as it does not support per Write TTL ...
		// FIXME nats has restrictions on the key, we cannot use slashes AFAICT
		// host, port, clusterid
		return natsjs.NewStore(
			append(storeopts,
				natsjs.NatsOptions(nats.Options{Name: "TODO"}),
				natsjs.DefaultTTL(options.TTL),
			)...,
		) // TODO test with ocis nats
	}
}

func storeOptions(o *Options) []store.Option {
	var opts []store.Option

	if o.Addresses != nil {
		opts = append(opts, store.Nodes(o.Addresses...))
	}

	if o.Database != "" {
		opts = append(opts, store.Database(o.Database))

	}

	if o.Table != "" {
		opts = append(opts, store.Table(o.Table))

	}

	return opts

}
