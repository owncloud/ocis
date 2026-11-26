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
	"bytes"
	"container/list"
	"fmt"
	"sync"

	"github.com/owncloud/reva/v2/pkg/store"
	"github.com/shamaton/msgpack/v2"
	microstore "go-micro.dev/v4/store"
)

const (
	// defaultDecodedAttributesEntries bounds the number of decoded attribute
	// maps kept in the process-local decoded layer.
	defaultDecodedAttributesEntries = 16384
	// defaultDecodedAttributesBytes bounds the (estimated) memory used by the
	// process-local decoded layer.
	defaultDecodedAttributesBytes = 64 << 20
)

// fileMetadataCache caches node attributes in a (possibly shared) store.
//
// In addition to the store, it keeps a bounded, process-local layer of
// already decoded attribute maps (map[string][]byte). The store remains the
// single source of truth: every read still fetches the serialized value from
// the store and only skips the msgpack decode if the fetched bytes are
// byte-for-byte identical to the bytes the decoded entry was built from. A
// cache hit therefore returns exactly what decoding the current store value
// would return, regardless of which process wrote it.
//
// The decoded layer is nil (disabled) for noop stores and if it is disabled
// in the configuration.
type fileMetadataCache struct {
	cacheStore
	decoded *decodedAttributesCache
	decode  func([]byte, any) error
}

// NewFileMetadataCache creates a new FileMetadataCache
func NewFileMetadataCache(cfg Config) FileMetadataCache {
	c := &fileMetadataCache{
		decoded: newDecodedAttributesCacheFromConfig(cfg),
		decode:  msgpack.Unmarshal,
	}
	c.s = getStore(cfg)
	c.database = cfg.Database
	c.table = cfg.Table
	c.ttl = cfg.TTL

	return c
}

// newDecodedAttributesCacheFromConfig returns the decoded layer for cfg, or nil
// if it is disabled: for noop stores, which never return a value, and if one
// of the limits is negative. A limit of 0 uses the built-in default.
func newDecodedAttributesCacheFromConfig(cfg Config) *decodedAttributesCache {
	if cfg.Store == store.TypeNoop || cfg.DecodedAttributesMaxEntries < 0 || cfg.DecodedAttributesMaxBytes < 0 {
		return nil
	}
	maxEntries, maxBytes := cfg.DecodedAttributesMaxEntries, cfg.DecodedAttributesMaxBytes
	if maxEntries == 0 {
		maxEntries = defaultDecodedAttributesEntries
	}
	if maxBytes == 0 {
		maxBytes = defaultDecodedAttributesBytes
	}
	return newDecodedAttributesCache(maxEntries, maxBytes)
}

// PullFromCache reads the value for the given key from the store. Attribute
// maps (*map[string][]byte) are served from the decoded layer if the store
// still holds the exact bytes they were decoded from.
func (c *fileMetadataCache) PullFromCache(key string, dest interface{}) error {
	attrs, ok := dest.(*map[string][]byte)
	if !ok || attrs == nil || c.decoded == nil {
		return c.cacheStore.PullFromCache(key, dest)
	}

	r, err := c.s.Read(key, microstore.ReadFrom(c.database, c.table), microstore.ReadLimit(1))
	if err != nil {
		c.decoded.remove(key)
		return err
	}
	if len(r) == 0 {
		c.decoded.remove(key)
		return fmt.Errorf("not found")
	}
	raw := r[0].Value

	if cached, ok := c.decoded.get(key, raw); ok {
		*attrs = copyAttributes(*attrs, cached)
		return nil
	}

	// Keep a private copy of the raw bytes: the decoded values alias them and
	// they are used to validate later reads.
	owned := bytes.Clone(raw)
	decoded := map[string][]byte{}
	if err := c.decode(owned, &decoded); err != nil {
		return err
	}
	c.decoded.add(key, owned, decoded, true)
	*attrs = copyAttributes(*attrs, decoded)
	return nil
}

// PushToCache writes the value for the given key to the store. Attribute maps
// are also put into the decoded layer so that the next read does not have to
// decode them.
func (c *fileMetadataCache) PushToCache(key string, src interface{}) error {
	if c.decoded == nil {
		return c.cacheStore.PushToCache(key, src)
	}
	attrs, ok := src.(map[string][]byte)
	if !ok {
		c.decoded.remove(key)
		return c.cacheStore.PushToCache(key, src)
	}

	b, err := msgpack.Marshal(attrs)
	if err != nil {
		c.decoded.remove(key)
		return err
	}
	err = c.s.Write(
		&microstore.Record{Key: key, Value: b, Expiry: c.ttl},
		microstore.WriteTo(c.database, c.table),
		microstore.WriteTTL(c.ttl),
	)
	if err != nil {
		c.decoded.remove(key)
		return err
	}
	// copy the source, the caller still owns it and might modify it
	c.decoded.add(key, b, copyAttributes(nil, attrs), false)
	return nil
}

// Delete deletes the given key from the store and the decoded layer
func (c *fileMetadataCache) Delete(key string, opts ...microstore.DeleteOption) error {
	c.decoded.remove(key)
	return c.cacheStore.Delete(key, opts...)
}

// RemoveMetadata removes a reference from the metadata cache
func (c *fileMetadataCache) RemoveMetadata(path string) error {
	return c.Delete(path)
}

// copyAttributes copies all attributes from src into dst and returns dst. If
// dst is empty a new, presized map is allocated. All values are copied into a
// single new backing array, each value is capped to its length so that
// appending to it in dst can never write into another value.
func copyAttributes(dst, src map[string][]byte) map[string][]byte {
	if len(dst) == 0 {
		dst = make(map[string][]byte, len(src))
	}
	total := 0
	for _, v := range src {
		total += len(v)
	}
	arena := make([]byte, total)
	for k, v := range src {
		if v == nil {
			dst[k] = nil
			continue
		}
		n := copy(arena, v)
		dst[k] = arena[:n:n]
		arena = arena[n:]
	}
	return dst
}

// decodedAttributesCache is a bounded LRU of decoded attribute maps. Entries
// are immutable once added and are only ever handed out as copies.
type decodedAttributesCache struct {
	mu         sync.Mutex
	maxEntries int
	maxBytes   int
	bytes      int
	ll         *list.List
	entries    map[string]*list.Element
}

type decodedAttributesEntry struct {
	key   string
	raw   []byte
	attrs map[string][]byte
	cost  int
}

func newDecodedAttributesCache(maxEntries, maxBytes int) *decodedAttributesCache {
	return &decodedAttributesCache{
		maxEntries: maxEntries,
		maxBytes:   maxBytes,
		ll:         list.New(),
		entries:    map[string]*list.Element{},
	}
}

// The overhead constants are tuned for 64-bit Go maps. If
// TestFileMetadataCacheDecodedLayerAccounting drifts after a Go upgrade,
// retune them rather than widening the test's tolerance.
const (
	// entryOverhead approximates the fixed heap cost of an entry: list
	// element, entry struct, attribute map header and the index map slot.
	entryOverhead = 256
	// attributeOverhead approximates the per attribute heap cost of a
	// map[string][]byte: string and slice headers, control byte, load factor
	// slack and allocation rounding of the key.
	attributeOverhead = 64
)

// entryCost estimates the heap retained by an entry. If valuesAliasRaw is
// true the attribute values are sub slices of raw (as returned by the msgpack
// decoder) and must not be counted again.
func entryCost(key string, raw []byte, attrs map[string][]byte, valuesAliasRaw bool) int {
	cost := entryOverhead + len(key) + cap(raw)
	for k, v := range attrs {
		cost += attributeOverhead + len(k)
		if !valuesAliasRaw {
			cost += len(v)
		}
	}
	return cost
}

// get returns the decoded attributes for key if they were decoded from raw.
// The returned map must not be modified.
//
// The comparison with raw runs without holding the lock: entries are
// immutable once added, so only the lookup and the LRU bookkeeping need it.
// The bookkeeping is skipped if the key was replaced or removed in between.
func (c *decodedAttributesCache) get(key string, raw []byte) (map[string][]byte, bool) {
	c.mu.Lock()
	el, ok := c.entries[key]
	var e *decodedAttributesEntry
	if ok {
		e = el.Value.(*decodedAttributesEntry)
	}
	c.mu.Unlock()
	if !ok {
		return nil, false
	}

	hit := bytes.Equal(e.raw, raw)

	c.mu.Lock()
	defer c.mu.Unlock()
	if c.entries[key] == el {
		if hit {
			c.ll.MoveToFront(el)
		} else {
			// the store holds a different value, this entry is stale
			c.removeElement(el)
		}
	}
	if !hit {
		return nil, false
	}
	return e.attrs, true
}

// add stores the decoded attributes for key. Neither raw nor attrs may be
// modified afterwards. valuesAliasRaw tells whether the values of attrs are
// sub slices of raw.
func (c *decodedAttributesCache) add(key string, raw []byte, attrs map[string][]byte, valuesAliasRaw bool) {
	e := &decodedAttributesEntry{key: key, raw: raw, attrs: attrs, cost: entryCost(key, raw, attrs, valuesAliasRaw)}
	c.mu.Lock()
	defer c.mu.Unlock()
	if el, ok := c.entries[key]; ok {
		c.removeElement(el)
	}
	if e.cost > c.maxBytes {
		return
	}
	c.entries[key] = c.ll.PushFront(e)
	c.bytes += e.cost
	for c.ll.Len() > c.maxEntries || c.bytes > c.maxBytes {
		c.removeElement(c.ll.Back())
	}
}

func (c *decodedAttributesCache) remove(key string) {
	if c == nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if el, ok := c.entries[key]; ok {
		c.removeElement(el)
	}
}

func (c *decodedAttributesCache) removeElement(el *list.Element) {
	e := c.ll.Remove(el).(*decodedAttributesEntry)
	delete(c.entries, e.key)
	c.bytes -= e.cost
}

func (c *decodedAttributesCache) len() int {
	if c == nil {
		return 0
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.ll.Len()
}

func (c *decodedAttributesCache) size() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.bytes
}
