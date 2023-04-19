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
	"strings"
	"time"
)

type fileMetadataCache struct {
	cacheStore
}

// NewFileMetadataCache creates a new FileMetadataCache
func NewFileMetadataCache(store string, nodes []string, database, table string, ttl time.Duration) FileMetadataCache {
	c := &fileMetadataCache{}
	c.s = getStore(store, nodes, database, table, ttl)
	c.database = database
	c.table = table
	c.ttl = ttl

	return c
}

// RemoveMetadata removes a reference from the metadata cache
func (c *fileMetadataCache) RemoveMetadata(path string) error {
	keys, err := c.List()
	if err != nil {
		return err
	}

	for _, key := range keys {
		if strings.HasPrefix(key, path) {
			_ = c.Delete(key)
		}
	}
	return nil
}
