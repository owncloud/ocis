// Copyright 2018-2023 CERN
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

package store

import (
	"context"
	"strings"
	"time"

	natsjs "github.com/go-micro/plugins/v4/store/nats-js"
	natsjskv "github.com/go-micro/plugins/v4/store/nats-js-kv"
	"github.com/go-micro/plugins/v4/store/redis"
	redisopts "github.com/go-redis/redis/v8"
	"github.com/nats-io/nats.go"
	"github.com/owncloud/reva/v2/pkg/store/etcd"
	"github.com/owncloud/reva/v2/pkg/store/memory"
	"go-micro.dev/v4/logger"
	microstore "go-micro.dev/v4/store"
)

var ocMemStore *microstore.Store

const (
	// TypeMemory represents memory stores
	TypeMemory = "memory"
	// TypeNoop represents noop stores
	TypeNoop = "noop"
	// TypeEtcd represents etcd stores
	TypeEtcd = "etcd"
	// TypeRedis represents redis stores
	TypeRedis = "redis"
	// TypeRedisSentinel represents redis-sentinel stores
	TypeRedisSentinel = "redis-sentinel"
	// TypeOCMem represents ocmem stores
	TypeOCMem = "ocmem"
	// TypeNatsJS represents nats-js stores
	TypeNatsJS = "nats-js"
	// TypeNatsJSKV represents nats-js-kv stores
	TypeNatsJSKV = "nats-js-kv"
)

// Create initializes a new store
func Create(opts ...microstore.Option) microstore.Store {
	options := &microstore.Options{
		Context: context.Background(),
	}
	for _, o := range opts {
		o(options)
	}

	storeType, _ := options.Context.Value(typeContextKey{}).(string)

	switch storeType {
	case TypeNoop:
		return microstore.NewNoopStore(opts...)
	case TypeEtcd:
		return etcd.NewStore(opts...)
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
			microstore.Database(options.Database),
			microstore.Table(options.Table),
			microstore.Nodes(redisNodes...),
			redis.WithRedisOptions(redisopts.UniversalOptions{
				MasterName: redisMaster,
			}),
		)
	case TypeOCMem:
		if ocMemStore == nil {
			var memStore microstore.Store

			sizeNum, _ := options.Context.Value(sizeContextKey{}).(int)
			if sizeNum <= 0 {
				memStore = memory.NewMultiMemStore()
			} else {
				memStore = memory.NewMultiMemStore(
					microstore.WithContext(
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
		ttl, _ := options.Context.Value(ttlContextKey{}).(time.Duration)
		if mem, _ := options.Context.Value(disablePersistanceContextKey{}).(bool); mem {
			opts = append(opts, natsjs.DefaultMemory())
		}
		// TODO nats needs a DefaultTTL option as it does not support per Write TTL ...
		// FIXME nats has restrictions on the key, we cannot use slashes AFAICT
		// host, port, clusterid
		natsOptions := defaultNatsOptions(options)
		return natsjs.NewStore(
			append(opts,
				natsjs.NatsOptions(natsOptions), // always pass in properly initialized default nats options
				natsjs.DefaultTTL(ttl))...,
		) // TODO test with ocis nats
	case TypeNatsJSKV:
		// NOTE: nats needs a DefaultTTL option as it does not support per Write TTL ...
		ttl, _ := options.Context.Value(ttlContextKey{}).(time.Duration)
		if mem, _ := options.Context.Value(disablePersistanceContextKey{}).(bool); mem {
			opts = append(opts, natsjskv.DefaultMemory())
		}

		natsOptions := defaultNatsOptions(options)
		return natsjskv.NewStore(
			append(opts,
				natsjskv.NatsOptions(natsOptions), // always pass in properly initialized default nats options
				natsjskv.EncodeKeys(),
				natsjskv.DefaultTTL(ttl))...,
		)
	case TypeMemory, "mem", "": // allow existing short form and use as default
		return microstore.NewMemoryStore(opts...)
	default:
		// try to log an error
		if options.Logger == nil {
			options.Logger = logger.DefaultLogger
		}
		options.Logger.Logf(logger.ErrorLevel, "unknown store type: '%s', falling back to memory", storeType)
		return microstore.NewMemoryStore(opts...)
	}
}

// defaultNatsOptions builds the nats.Options shared by the nats-js and
// nats-js-kv store backends. It overrides the NATS client defaults so the
// client never permanently gives up on a closed connection: the default
// MaxReconnect (60) combined with the default ReconnectWait (2s) means any
// NATS outage longer than ~2 minutes leaves the client permanently closed,
// which the store plugins then surface as "nats: connection closed" on every
// subsequent operation. Reconnecting forever, together with the connection
// state handlers, keeps the client alive and makes the transitions visible.
func defaultNatsOptions(options *microstore.Options) nats.Options {
	natsOptions := nats.GetDefaultOptions()
	natsOptions.Name = "reva-store"
	natsOptions.MaxReconnect = -1 // reconnect forever; the default of 60 gives up after ~2 minutes
	natsOptions.ReconnectWait = 5 * time.Second
	if auth, ok := options.Context.Value(authenticationContextKey{}).([]string); ok && len(auth) == 2 {
		natsOptions.User = auth[0]
		natsOptions.Password = auth[1]
	}
	natsOptions.DisconnectedErrCB = func(_ *nats.Conn, err error) {
		logger.Logf(logger.WarnLevel, "reva-store: nats connection disconnected: %v", err)
	}
	natsOptions.ReconnectedCB = func(c *nats.Conn) {
		logger.Logf(logger.InfoLevel, "reva-store: nats connection reconnected to %s", c.ConnectedUrl())
	}
	natsOptions.ClosedCB = func(_ *nats.Conn) {
		logger.Logf(logger.ErrorLevel, "reva-store: nats connection closed")
	}
	return natsOptions
}
