package memory

import (
	"context"

	"github.com/owncloud/ocis/v2/ocis-pkg/cache/memory"
	"go-micro.dev/v4/store"
)

// Prepare a context to be used with the memory implementation. The context
// is used to set up custom parameters to the specific implementation.
// In this case, you can configure the maximum capacity for the MemStore
// implementation as shown below.
// ```
// cache := NewMemStore(
//
//	store.WithContext(
//	  NewContext(
//	    ctx,
//	    map[string]interface{}{
//	      "maxCap": 50,
//	    },
//	  ),
//	),
//
// )
// ```
//
// Available options for the MemStore are:
// * "maxCap" -> 512 (int) The maximum number of elements the cache will hold.
// Adding additional elements will remove old elements to ensure we aren't over
// the maximum capacity.
//
// For convenience, this can also be used for the MultiMemStore.
// Deprecated: use "github.com/owncloud/ocis/v2/ocis-pkg/store/memory" NewContext
func NewContext(ctx context.Context, storeParams map[string]interface{}) context.Context {
	return memory.NewContext(ctx, storeParams)
}

// Create a new MemStore instance
// Deprecated: use "github.com/owncloud/ocis/v2/ocis-pkg/store/memory" NewMemStore
func NewMemStore(opts ...store.Option) store.Store {
	return memory.NewMemStore(opts...)
}
