package cache

import (
	"github.com/owncloud/ocis/v2/ocis-pkg/cache"
	"go-micro.dev/v4/store"
)

// Create returns a configured key-value micro store
//
// Each microservice (or whatever piece is using the store) should use the
// options available in the interface's operations to choose the right database
// and table to prevent collisions with other microservices.
// Recommended approach is to use "services" or "ocis-pkg" for the database,
// and "services/<service-name>/" or "ocis-pkg/<pkg>/" for the package name.
// Deprecated: use "github.com/owncloud/ocis/v2/ocis-pkg/store" Create
func Create(opts ...store.Option) store.Store {
	return cache.Create(opts...)
}
