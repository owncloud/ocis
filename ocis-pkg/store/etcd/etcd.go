package etcd

import (
	"github.com/owncloud/ocis/v2/ocis-pkg/cache/etcd"
	"go-micro.dev/v4/store"
)

// Create a new go-micro store backed by etcd
// Deprecated: use "github.com/owncloud/ocis/v2/ocis-pkg/store/etcd" NewEtcdStore
func NewEtcdStore(opts ...store.Option) store.Store {
	return etcd.NewEtcdStore(opts...)
}
