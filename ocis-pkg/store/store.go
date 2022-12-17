package store

import (
	"context"
	"strings"

	"github.com/owncloud/ocis/v2/ocis-pkg/store/etcd"
	"github.com/owncloud/ocis/v2/ocis-pkg/store/memory"
	"go-micro.dev/v4/store"
)

var ocMemStore *store.Store

// OcisStoreOptions
// FIXME: nolint
// nolint: revive
type OcisStoreOptions struct {
	Type    string
	Address string
	Size    int
}

// GetStore gets the configured key-value store to be used.
//
// Each microservice (or whatever piece is using the store) should use the
// options available in the interface's operations to choose the right database
// and table to prevent collisions with other microservices.
// Recommended approach is to use "services" or "ocis-pkg" for the database,
// and "services/<service-name>/" or "ocis-pkg/<pkg>/" for the package name.
//
// So far, only the name of the store and the node addresses are configurable
// via environment variables.
// Available options for "OCIS_STORE" are:
// * "noop", for a noop store (it does nothing)
// * "etcd", for etcd
// * "ocmem", custom in-memory implementation, with fixed size and optimized prefix
// and suffix search
// * "memory", for a in-memory implementation, which is the default if noone matches
//
// "OCIS_STORE_ADDRESS" is a comma-separated list of nodes that the store
// will use. This is currently usable only with the etcd implementation. If it
// isn't provided, "127.0.0.1:2379" will be the only node used.
//
// "OCIS_STORE_OCMEM_SIZE" will configure the maximum capacity of the cache for
// the "ocmem" implementation, in number of items that the cache can hold per table.
// You can use "OCIS_STORE_OCMEM_SIZE=5000" so the cache will hold up to 5000 elements.
// The parameter only affects to the "ocmem" implementation, the rest will ignore it.
// If an invalid value is used, the default will be used instead, so up to 512 elements
// the cache will hold.
func GetStore(ocisOpts OcisStoreOptions) store.Store {
	var s store.Store

	addresses := strings.Split(ocisOpts.Address, ",")
	opts := []store.Option{
		store.Nodes(addresses...),
	}

	switch ocisOpts.Type {
	case "noop":
		s = store.NewNoopStore(opts...)
	case "etcd":
		s = etcd.NewEtcdStore(opts...)
	case "ocmem":
		if ocMemStore == nil {
			var memStore store.Store

			sizeNum := ocisOpts.Size
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
		s = *ocMemStore
	default:
		s = store.NewMemoryStore(opts...)
	}
	return s
}
