//  Copyright (c) 2026 Couchbase, Inc.
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
	"math"

	"github.com/RoaringBitmap/roaring/v2"
	index "github.com/blevesearch/bleve_index_api"
)

type nestedIndexCache struct {
	cache *nestedCacheEntry
}

// newNestedIndexCache creates a new nested index cache
// instance, which contains cached edge list
// for a nested segment
func newNestedIndexCache() *nestedIndexCache {
	return &nestedIndexCache{}
}

// Clear clears the nested index cache, removing the cached edge list
func (nc *nestedIndexCache) Clear() {
	nc.cache = nil
}

func (nc *nestedIndexCache) initialize(numDocs uint64, edgeListOffset uint64, mem []byte) error {
	// pos stores the current read position
	pos := edgeListOffset
	if pos == 0 {
		// no edge list
		return nil
	}
	// read number of edges in the edge list
	numEdges, read := binary.Uvarint(mem[pos : pos+binary.MaxVarintLen64])
	if read <= 0 {
		return fmt.Errorf("error reading number of edges in nested edge list")
	}
	pos += uint64(read)
	// if no documents or edges/nested documents, return
	if numDocs == 0 || numEdges == 0 {
		return nil
	}
	edgeList := NewEdgeList(numDocs, numEdges)
	for i := uint64(0); i < numEdges; i++ {
		child, read := binary.Uvarint(mem[pos : pos+binary.MaxVarintLen64])
		if read <= 0 {
			return fmt.Errorf("error reading child doc id in nested edge list")
		}
		pos += uint64(read)
		parent, read := binary.Uvarint(mem[pos : pos+binary.MaxVarintLen64])
		if read <= 0 {
			return fmt.Errorf("error reading parent doc id in nested edge list")
		}
		pos += uint64(read)
		edgeList.AddEdge(child, parent)
	}
	nc.cache = &nestedCacheEntry{
		el: edgeList,
	}
	return nil
}

type nestedCacheEntry struct {
	// edgeList[child] = parent
	el EdgeList
}

func (nc *nestedIndexCache) ancestry(docNum uint64, prealloc []index.AncestorID) []index.AncestorID {
	cache := nc.cache
	// add self as first ancestor
	prealloc = append(prealloc, index.NewAncestorID(docNum))
	if cache == nil || cache.el == nil {
		return prealloc
	}
	current := docNum
	for {
		parent, ok := cache.el.Parent(current)
		if !ok {
			break
		}
		prealloc = append(prealloc, index.NewAncestorID(parent))
		current = parent
	}
	return prealloc
}

func (nc *nestedIndexCache) edgeList() EdgeList {
	cache := nc.cache
	if cache == nil || cache.el == nil {
		return nil
	}
	return cache.el
}

func (nc *nestedIndexCache) countNested() uint64 {
	cache := nc.cache
	if cache == nil || cache.el == nil {
		return 0
	}
	return cache.el.Count()
}

// countRoot returns the number of root documents in the given bitmap
func (nc *nestedIndexCache) countRoot(bm *roaring.Bitmap) uint64 {
	var totalDocs uint64
	if bm == nil {
		// if bitmap is empty, return 0
		return totalDocs
	}
	totalDocs = bm.GetCardinality()
	cache := nc.cache
	if cache == nil || cache.el == nil {
		// if cache is nil, no nested docs, so all docs are root docs
		// so just return the cardinality of the bitmap
		return totalDocs
	}
	// count nested documents in the bitmap, a nested doc is one that has a parent in the edge list
	var nestedDocCount uint64
	bm.Iterate(func(docNum uint32) bool {
		if _, ok := cache.el.Parent(uint64(docNum)); ok {
			nestedDocCount++
		}
		return true
	})
	// root docs = total docs - nested docs
	if totalDocs < nestedDocCount {
		// should not happen, but just in case
		return 0
	}
	return totalDocs - nestedDocCount
}

// -------------------------------------------------------

// EdgeList provides an interface to access parent of a child document
type EdgeList interface {
	// Parent returns the parent of the given child document ID,
	// and a boolean indicating if the parent exists.
	Parent(child uint64) (uint64, bool)

	// AddEdge adds an edge from child to parent in the edge list.
	AddEdge(child uint64, parent uint64)

	// Count returns the number of edges in the edge list.
	Count() uint64

	// Iterate iterates over all edges in the edge list, calling the provided function
	// with each child-parent pair. If the function returns false, iteration stops.
	Iterate(func(child uint64, parent uint64) bool)
}

type edgeListMap struct {
	edges map[uint64]uint64
}

func newEdgeListMap(numEdges uint64) *edgeListMap {
	return &edgeListMap{
		edges: make(map[uint64]uint64, numEdges),
	}
}

func (elm *edgeListMap) Parent(child uint64) (uint64, bool) {
	parent, ok := elm.edges[child]
	return parent, ok
}

func (elm *edgeListMap) AddEdge(child uint64, parent uint64) {
	elm.edges[child] = parent
}

func (elm *edgeListMap) Count() uint64 {
	return uint64(len(elm.edges))
}

func (elm *edgeListMap) Iterate(f func(child uint64, parent uint64) bool) {
	for child, parent := range elm.edges {
		if !f(child, parent) {
			return
		}
	}
}

type edgeListSlice struct {
	count    uint64
	sentinel uint64
	edges    []uint64
}

func newEdgeListSlice(numDocs uint64, numEdges uint64) *edgeListSlice {
	var sentinel uint64 = math.MaxUint64
	edges := make([]uint64, numDocs)
	for i := range edges {
		edges[i] = sentinel
	}
	return &edgeListSlice{
		count:    numEdges,
		sentinel: sentinel,
		edges:    edges,
	}
}

func (els *edgeListSlice) Parent(child uint64) (uint64, bool) {
	if child >= uint64(len(els.edges)) {
		return 0, false
	}
	parent := els.edges[child]
	if parent == els.sentinel {
		return 0, false
	}
	return parent, true
}

func (el *edgeListSlice) AddEdge(child uint64, parent uint64) {
	if child >= uint64(len(el.edges)) {
		// out of bounds, ignore as this should not happen
		return
	}
	el.edges[child] = parent
}

func (el *edgeListSlice) Count() uint64 {
	return el.count
}

func (el *edgeListSlice) Iterate(f func(child uint64, parent uint64) bool) {
	for child, parent := range el.edges {
		if parent != el.sentinel {
			if !f(uint64(child), parent) {
				return
			}
		}
	}
}

// nestedCacheRatio defines the threshold ratio of nested documents to total documents.
// It is derived using the following reasoning:
//
// Let N = number of nested documents (i.e., edges in the edge list)
// Let T = total number of documents
//
// Memory usage if the edge list is stored as a map[uint64]uint64:
//
//	~30 bytes per entry (key + value + map overhead)
//	Total ≈ 30 * N bytes
//
// Memory usage if the edge list is stored as a []uint64:
//
//	8 bytes per entry
//	Total ≈ 8 * T bytes
//
// We want the threshold at which a map becomes more memory-efficient than a slice:
//
//	30N < 8T
//	N/T < 8/30
//
// Therefore, if the ratio of nested documents to total documents is less than 8/30,
// we use a map for the edge list; otherwise, we use a slice.
var edgeListMapThreshold = 8.0 / 30.0

// NewEdgeList creates a new EdgeList instance based on the provided
// constants, the total number of documents and the number of nested documents/edges.
func NewEdgeList(numDocs uint64, numEdges uint64) EdgeList {
	if numDocs == 0 || numEdges == 0 {
		// no edges, return nil
		return nil
	}
	ratio := float64(numEdges) / float64(numDocs)
	if ratio < edgeListMapThreshold {
		// use map representation
		return newEdgeListMap(numEdges)
	}
	// use slice representation
	return newEdgeListSlice(numDocs, numEdges)
}
