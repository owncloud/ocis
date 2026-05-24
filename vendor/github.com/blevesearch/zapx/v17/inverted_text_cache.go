//  Copyright (c) 2025 Couchbase, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 		http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package zap

import (
	"encoding/binary"
	"fmt"
	"sync"

	"github.com/blevesearch/vellum"
)

func newInvertedIndexCache() *invertedIndexCache {
	return &invertedIndexCache{
		cache: make(map[uint16]*invertedCacheEntry),
	}
}

type invertedIndexCache struct {
	m sync.RWMutex

	cache map[uint16]*invertedCacheEntry
}

func (sc *invertedIndexCache) Clear() {
	sc.m.Lock()
	sc.cache = nil
	sc.m.Unlock()
}

// loadOrCreate loads the inverted index cache for the specified fieldID if it is already present,
// or creates it if not. The inverted index cache for a fieldID consists of an FST (Finite State Transducer):
// - A Vellum FST (Finite State Transducer) representing the TermDictionary.
// This function returns the loaded or newly created FST, and the number of bytes read from the provided memory slice,
// if the cache was created.
func (sc *invertedIndexCache) loadOrCreate(fieldID uint16, mem []byte, fr *FileReader) (*vellum.FST, uint64, error) {
	sc.m.RLock()
	entry, ok := sc.cache[fieldID]
	if ok {
		sc.m.RUnlock()
		return entry.load()
	}

	sc.m.RUnlock()

	sc.m.Lock()
	defer sc.m.Unlock()

	entry, ok = sc.cache[fieldID]
	if ok {
		return entry.load()
	}

	return sc.createAndCacheLOCKED(fieldID, mem, fr)
}

// createAndCacheLOCKED creates the inverted index cache for the specified fieldID and caches it.
func (sc *invertedIndexCache) createAndCacheLOCKED(fieldID uint16, mem []byte, fr *FileReader) (*vellum.FST, uint64, error) {
	var pos uint64
	vellumLen, read := binary.Uvarint(mem[pos : pos+binary.MaxVarintLen64])
	if vellumLen == 0 || read <= 0 {
		return nil, 0, fmt.Errorf("vellum length is 0")
	}
	pos += uint64(read)
	fstBytes, err := fr.process(mem[pos : pos+vellumLen])
	if err != nil {
		return nil, 0, fmt.Errorf("error processing vellum bytes: %v", err)
	}
	fst, err := vellum.Load(fstBytes)
	if err != nil {
		return nil, 0, fmt.Errorf("vellum err: %v", err)
	}
	pos += vellumLen
	sc.insertLOCKED(fieldID, fst)
	return fst, pos, nil
}

// insertLOCKED inserts the vellum FST into the cache for the specified fieldID.
func (sc *invertedIndexCache) insertLOCKED(fieldID uint16, fst *vellum.FST) {
	_, ok := sc.cache[fieldID]
	if !ok {
		sc.cache[fieldID] = &invertedCacheEntry{
			fst: fst,
		}
	}
}

// invertedCacheEntry is the vellum FST and is the value stored in the invertedIndexCache cache, for a given fieldID.
type invertedCacheEntry struct {
	fst *vellum.FST
}

func (ce *invertedCacheEntry) load() (*vellum.FST, uint64, error) {
	return ce.fst, 0, nil
}
