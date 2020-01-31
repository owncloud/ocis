// Package registry provides accessors to runtime services
package registry

import (
	"sync"

	mstore "github.com/micro/go-micro/v2/store"
	store "github.com/owncloud/ocis-accounts/pkg/store/filesystem"
)

var (
	once *sync.Once = &sync.Once{}
	// Store is a micro store implementation
	Store mstore.Store
)

func init() {
	once.Do(func() {
		Store = store.New()
	})
}
