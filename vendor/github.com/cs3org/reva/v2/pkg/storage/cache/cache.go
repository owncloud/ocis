// Copyright 2018-2021 CERN
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// In applying this license, CERN does not waive the privileges and immunities
// granted to it by virtue of its status as an Intergovernmental Organization
// or submit itself to any jurisdiction.

package cache

import (
	"fmt"
	"strings"
	"sync"
	"time"

	userpb "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	natsjs "github.com/go-micro/plugins/v4/store/nats-js"
	"github.com/go-micro/plugins/v4/store/redis"
	redisopts "github.com/go-redis/redis/v8"
	"github.com/nats-io/nats.go"
	microetcd "github.com/owncloud/ocis/v2/ocis-pkg/store/etcd"
	"github.com/shamaton/msgpack/v2"
	microstore "go-micro.dev/v4/store"
)

var (
	// DefaultStatCache is the memory store.
	statCaches                = make(map[string]StatCache)
	providerCaches            = make(map[string]ProviderCache)
	createHomeCaches          = make(map[string]CreateHomeCache)
	createPersonalSpaceCaches = make(map[string]CreatePersonalSpaceCache)
	fileMetadataCaches        = make(map[string]FileMetadataCache)
	mutex                     sync.Mutex
)

// Cache handles key value operations on caches
type Cache interface {
	PullFromCache(key string, dest interface{}) error
	PushToCache(key string, src interface{}) error
	List(opts ...microstore.ListOption) ([]string, error)
	Delete(key string, opts ...microstore.DeleteOption) error
	Close() error
}

// StatCache handles removing keys from a stat cache
type StatCache interface {
	Cache
	RemoveStat(userID *userpb.UserId, res *provider.ResourceId)
	GetKey(userID *userpb.UserId, ref *provider.Reference, metaDataKeys, fieldMaskPaths []string) string
}

// ProviderCache handles removing keys from a storage provider cache
type ProviderCache interface {
	Cache
	RemoveListStorageProviders(res *provider.ResourceId)
	GetKey(userID *userpb.UserId, spaceID string) string
}

// CreateHomeCache handles removing keys from a create home cache
type CreateHomeCache interface {
	Cache
	RemoveCreateHome(res *provider.ResourceId)
	GetKey(userID *userpb.UserId) string
}

// CreatePersonalSpaceCache handles removing keys from a create home cache
type CreatePersonalSpaceCache interface {
	Cache
	GetKey(userID *userpb.UserId) string
}

// FileMetadataCache handles file metadata
type FileMetadataCache interface {
	Cache
	RemoveMetadata(path string) error
}

// GetStatCache will return an existing StatCache for the given store, nodes, database and table
// If it does not exist yet it will be created, different TTLs are ignored
func GetStatCache(cacheStore string, cacheNodes []string, database, table string, ttl time.Duration) StatCache {
	mutex.Lock()
	defer mutex.Unlock()

	key := strings.Join(append(append([]string{cacheStore}, cacheNodes...), database, table), ":")
	if statCaches[key] == nil {
		statCaches[key] = NewStatCache(cacheStore, cacheNodes, database, table, ttl)
	}
	return statCaches[key]
}

// GetProviderCache will return an existing ProviderCache for the given store, nodes, database and table
// If it does not exist yet it will be created, different TTLs are ignored
func GetProviderCache(cacheStore string, cacheNodes []string, database, table string, ttl time.Duration) ProviderCache {
	mutex.Lock()
	defer mutex.Unlock()

	key := strings.Join(append(append([]string{cacheStore}, cacheNodes...), database, table), ":")
	if providerCaches[key] == nil {
		providerCaches[key] = NewProviderCache(cacheStore, cacheNodes, database, table, ttl)
	}
	return providerCaches[key]
}

// GetCreateHomeCache will return an existing CreateHomeCache for the given store, nodes, database and table
// If it does not exist yet it will be created, different TTLs are ignored
func GetCreateHomeCache(cacheStore string, cacheNodes []string, database, table string, ttl time.Duration) CreateHomeCache {
	mutex.Lock()
	defer mutex.Unlock()

	key := strings.Join(append(append([]string{cacheStore}, cacheNodes...), database, table), ":")
	if createHomeCaches[key] == nil {
		createHomeCaches[key] = NewCreateHomeCache(cacheStore, cacheNodes, database, table, ttl)
	}
	return createHomeCaches[key]
}

// GetCreatePersonalSpaceCache will return an existing CreatePersonalSpaceCache for the given store, nodes, database and table
// If it does not exist yet it will be created, different TTLs are ignored
func GetCreatePersonalSpaceCache(cacheStore string, cacheNodes []string, database, table string, ttl time.Duration) CreatePersonalSpaceCache {
	mutex.Lock()
	defer mutex.Unlock()

	key := strings.Join(append(append([]string{cacheStore}, cacheNodes...), database, table), ":")
	if createPersonalSpaceCaches[key] == nil {
		createPersonalSpaceCaches[key] = NewCreatePersonalSpaceCache(cacheStore, cacheNodes, database, table, ttl)
	}
	return createPersonalSpaceCaches[key]
}

// GetFileMetadataCache will return an existing GetFileMetadataCache for the given store, nodes, database and table
// If it does not exist yet it will be created, different TTLs are ignored
func GetFileMetadataCache(cacheStore string, cacheNodes []string, database, table string, ttl time.Duration) FileMetadataCache {
	mutex.Lock()
	defer mutex.Unlock()

	key := strings.Join(append(append([]string{cacheStore}, cacheNodes...), database, table), ":")
	if fileMetadataCaches[key] == nil {
		fileMetadataCaches[key] = NewFileMetadataCache(cacheStore, cacheNodes, database, table, ttl)
	}
	return fileMetadataCaches[key]
}

// CacheStore holds cache store specific configuration
type cacheStore struct {
	s               microstore.Store
	database, table string
	ttl             time.Duration
}

// NewCache initializes a new CacheStore
func NewCache(store string, nodes []string, database, table string, ttl time.Duration) Cache {
	return cacheStore{
		s:        getStore(store, nodes, database, table, ttl), // some stores use a default ttl so we pass it when initializing
		database: database,
		table:    table,
		ttl:      ttl, // some stores use the ttl on every write, so we remember it here
	}
}

func getStore(store string, nodes []string, database, table string, ttl time.Duration) microstore.Store {
	switch store {
	case "etcd":
		return microetcd.NewEtcdStore(
			microstore.Nodes(nodes...),
			microstore.Database(database),
			microstore.Table(table),
		)
	case "nats-js":
		// TODO nats needs a DefaultTTL option as it does not support per Write TTL ...
		// FIXME nats has restrictions on the key, we cannot use slashes AFAICT
		// host, port, clusterid
		return natsjs.NewStore(
			microstore.Nodes(nodes...),
			microstore.Database(database),
			microstore.Table(table),
			natsjs.NatsOptions(nats.Options{Name: "TODO"}),
			natsjs.DefaultTTL(ttl),
		) // TODO test with ocis nats
	case "redis":
		return redis.NewStore(
			microstore.Database(database),
			microstore.Table(table),
			microstore.Nodes(nodes...),
		) // only the first node is taken into account
	case "redis-sentinel":
		redisMaster := ""
		redisNodes := []string{}
		for _, node := range nodes {
			parts := strings.SplitN(node, "/", 2)
			if len(parts) != 2 {
				return nil
			}
			// the first node is used to retrieve the redis master
			redisNodes = append(redisNodes, parts[0])
			if redisMaster == "" {
				redisMaster = parts[1]
			}
		}

		return redis.NewStore(
			microstore.Database(database),
			microstore.Table(table),
			microstore.Nodes(redisNodes...),
			redis.WithRedisOptions(redisopts.UniversalOptions{
				MasterName: redisMaster,
			}),
		)
	case "memory":
		return microstore.NewStore(
			microstore.Database(database),
			microstore.Table(table),
		)
	default:
		return microstore.NewNoopStore(
			microstore.Database(database),
			microstore.Table(table),
		)
	}
}

// PullFromCache pulls a value from the configured database and table of the underlying store using the given key
func (cache cacheStore) PullFromCache(key string, dest interface{}) error {
	r, err := cache.s.Read(key, microstore.ReadFrom(cache.database, cache.table), microstore.ReadLimit(1))
	if err != nil {
		return err
	}
	if len(r) == 0 {
		return fmt.Errorf("not found")
	}

	return msgpack.Unmarshal(r[0].Value, &dest)
}

// PushToCache pushes a key and value to the configured database and table of the underlying store
func (cache cacheStore) PushToCache(key string, src interface{}) error {
	b, err := msgpack.Marshal(src)
	if err != nil {
		return err
	}
	return cache.s.Write(
		&microstore.Record{Key: key, Value: b},
		microstore.WriteTo(cache.database, cache.table),
		microstore.WriteTTL(cache.ttl),
	)
}

// List lists the keys on the configured database and table of the underlying store
func (cache cacheStore) List(opts ...microstore.ListOption) ([]string, error) {
	o := []microstore.ListOption{
		microstore.ListFrom(cache.database, cache.table),
	}
	o = append(o, opts...)
	keys, err := cache.s.List(o...)
	if err != nil {
		return nil, err
	}
	for i, key := range keys {
		keys[i] = strings.TrimPrefix(key, cache.table)
	}
	return keys, nil
}

// Delete deletes the given key on the configured database and table of the underlying store
func (cache cacheStore) Delete(key string, opts ...microstore.DeleteOption) error {
	o := []microstore.DeleteOption{
		microstore.DeleteFrom(cache.database, cache.table),
	}
	o = append(o, opts...)
	return cache.s.Delete(key, o...)
}

// Close closes the underlying store
func (cache cacheStore) Close() error {
	return cache.s.Close()
}
