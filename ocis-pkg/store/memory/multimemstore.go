package memory

import (
	"github.com/owncloud/ocis/v2/ocis-pkg/cache/memory"
	"go-micro.dev/v4/store"
)

// Create a new MultiMemStore. A new MemStore will be mapped based on the options.
// A default MemStore will be mapped if no Database and Table aren't used.
// Deprecated: use "github.com/owncloud/ocis/v2/ocis-pkg/store/memory" NewMultiMemStore
func NewMultiMemStore(opts ...store.Option) store.Store {
	return memory.NewMultiMemStore(opts...)
}
