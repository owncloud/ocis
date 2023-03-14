package store

import (
	"context"
	"strings"

	natsjs "github.com/go-micro/plugins/v4/store/nats-js"
	"github.com/go-micro/plugins/v4/store/redis"
	redisopts "github.com/go-redis/redis/v8"
	"github.com/nats-io/nats.go"
	"github.com/owncloud/ocis/v2/ocis-pkg/store/etcd"
	"github.com/owncloud/ocis/v2/ocis-pkg/store/memory"
	"go-micro.dev/v4/store"
)

var ocMemStore *store.Store

const (
	TypeMem           = "mem"
	TypeNoop          = "noop"
	TypeEtcd          = "etcd"
	TypeRedis         = "redis"
	TypeRedisSentinel = "redis-sentinel"
	TypeOCMem         = "ocmem"
	TypeNatsJS        = "nats-js"
)

// Create returns a configured key-value store
//
// Each microservice (or whatever piece is using the store) should use the
// options available in the interface's operations to choose the right database
// and table to prevent collisions with other microservices.
// Recommended approach is to use "services" or "ocis-pkg" for the database,
// and "services/<service-name>/" or "ocis-pkg/<pkg>/" for the package name.
func Create(opts ...store.Option) store.Store {
	options := &store.Options{}
	for _, o := range opts {
		o(options)
	}

	cacheOpts, ok := options.Context.Value(cacheOptionsContextKey{}).(CacheOptions)
	if !ok {
		cacheOpts = CacheOptions{}
	}

	switch cacheOpts.Type {
	default:
		// TODO: better to error in default case?
		fallthrough
	case TypeMem:
		return store.NewMemoryStore(opts...)
	case TypeNoop:
		return store.NewNoopStore(opts...)
	case TypeEtcd:
		return etcd.NewEtcdStore(opts...)
	case TypeRedis:
		// FIXME redis plugin does not support redis cluster or ring -> needs upstream patch or our implementation
		return redis.NewStore(opts...)
	case TypeRedisSentinel:
		redisMaster := ""
		redisNodes := []string{}
		for _, node := range options.Nodes {
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
			store.Database(options.Database),
			store.Table(options.Table),
			store.Nodes(redisNodes...),
			redis.WithRedisOptions(redisopts.UniversalOptions{
				MasterName: redisMaster,
			}),
		)
	case TypeOCMem:
		if ocMemStore == nil {
			var memStore store.Store

			sizeNum := cacheOpts.Size
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
	case TypeNatsJS:
		// TODO nats needs a DefaultTTL option as it does not support per Write TTL ...
		// FIXME nats has restrictions on the key, we cannot use slashes AFAICT
		// host, port, clusterid
		return natsjs.NewStore(
			append(opts,
				natsjs.NatsOptions(nats.Options{Name: "TODO"}),
				natsjs.DefaultTTL(cacheOpts.TTL))...,
		) // TODO test with ocis nats
	}
}
