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

package lookup

import (
	"context"

	microstore "go-micro.dev/v4/store"

	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/options"
	"github.com/cs3org/reva/v2/pkg/store"
)

type StoreIDCache struct {
	cache microstore.Store
}

// NewMemoryIDCache returns a new MemoryIDCache
func NewStoreIDCache(o *options.Options) *StoreIDCache {
	return &StoreIDCache{
		cache: store.Create(
			store.Store(o.IDCache.Store),
			store.TTL(o.IDCache.TTL),
			store.Size(o.IDCache.Size),
			microstore.Nodes(o.IDCache.Nodes...),
			microstore.Database(o.IDCache.Database),
			microstore.Table(o.IDCache.Table),
			store.DisablePersistence(o.IDCache.DisablePersistence),
			store.Authentication(o.IDCache.AuthUsername, o.IDCache.AuthPassword),
		),
	}
}

// Add adds a new entry to the cache
func (c *StoreIDCache) Set(_ context.Context, spaceID, nodeID, val string) error {
	return c.cache.Write(&microstore.Record{
		Key:   cacheKey(spaceID, nodeID),
		Value: []byte(val),
	})
}

// Get returns the value for a given key
func (c *StoreIDCache) Get(_ context.Context, spaceID, nodeID string) (string, bool) {
	records, err := c.cache.Read(cacheKey(spaceID, nodeID))
	if err != nil || len(records) == 0 {
		return "", false
	}
	return string(records[0].Value), true
}

func cacheKey(spaceid, nodeID string) string {
	return spaceid + "!" + nodeID
}
