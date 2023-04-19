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
	"time"

	"go-micro.dev/v4/store"
)

type typeContextKey struct{}

// Store determines the implementation:
//   - "memory", for a in-memory implementation, which is also the default if noone matches
//   - "noop", for a noop store (it does nothing)
//   - "etcd", for etcd
//   - "nats-js" for nats-js, needs to have TTL configured at creation
//   - "redis", for redis
//   - "redis-sentinel", for redis-sentinel
//   - "ocmem", custom in-memory implementation, with fixed size and optimized prefix
//     and suffix search
func Store(val string) store.Option {
	return func(o *store.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}

		o.Context = context.WithValue(o.Context, typeContextKey{}, val)
	}
}

type sizeContextKey struct{}

// Size configures the maximum capacity of the cache for the "ocmem" implementation,
// in number of items that the cache can hold per table.
// You can use 5000 to make the cache hold up to 5000 elements.
// The parameter only affects to the "ocmem" implementation, the rest will ignore it.
// If an invalid value is used, the default of 512 will be used instead.
func Size(val int) store.Option {
	return func(o *store.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}

		o.Context = context.WithValue(o.Context, sizeContextKey{}, val)
	}
}

type ttlContextKey struct{}

// TTL is the time to live for documents stored in the store
func TTL(val time.Duration) store.Option {
	return func(o *store.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}

		o.Context = context.WithValue(o.Context, ttlContextKey{}, val)
	}
}
