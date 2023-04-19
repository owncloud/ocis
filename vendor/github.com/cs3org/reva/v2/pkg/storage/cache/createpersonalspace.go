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

package cache

import (
	"time"

	userpb "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
)

// CreatePersonalSpaceCache can invalidate all create home related cache entries
type createPersonalSpaceCache struct {
	cacheStore
}

// NewCreatePersonalSpaceCache creates a new CreatePersonalSpaceCache
func NewCreatePersonalSpaceCache(store string, nodes []string, database, table string, ttl time.Duration, size int) CreatePersonalSpaceCache {
	c := &createPersonalSpaceCache{}
	c.s = getStore(store, nodes, database, table, ttl, size)
	c.database = database
	c.table = table
	c.ttl = ttl

	return c
}

func (c createPersonalSpaceCache) GetKey(userID *userpb.UserId) string {
	return userID.GetOpaqueId()
}
