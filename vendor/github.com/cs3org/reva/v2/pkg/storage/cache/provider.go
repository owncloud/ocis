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
	"strings"
	"time"

	userpb "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
)

// ProviderCache can invalidate all provider related cache entries
type providerCache struct {
	cacheStore
}

// NewProviderCache creates a new ProviderCache
func NewProviderCache(store string, nodes []string, database, table string, ttl time.Duration) ProviderCache {
	c := &providerCache{}
	c.s = getStore(store, nodes, database, table, ttl)
	c.database = database
	c.table = table
	c.ttl = ttl

	return c
}

// RemoveListStorageProviders removes a reference from the listproviders cache
func (c providerCache) RemoveListStorageProviders(res *provider.ResourceId) {
	if res == nil {
		return
	}
	sid := res.SpaceId

	keys, err := c.List()
	if err != nil {
		// FIXME log error
		return
	}
	// FIXME add context option to List, Read and Write to upstream
	for _, key := range keys {
		if strings.Contains(key, sid) {
			_ = c.Delete(key)
			continue
		}
	}
}

func (c providerCache) GetKey(userID *userpb.UserId, spaceID string) string {
	if key := userID.GetOpaqueId() + "!" + spaceID; key != "!" {
		return key
	}
	return ""
}
