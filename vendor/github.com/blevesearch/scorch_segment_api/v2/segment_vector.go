//  Copyright (c) 2023 Couchbase, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build vectors
// +build vectors

package segment

import (
	"encoding/json"

	"github.com/RoaringBitmap/roaring/v2"
	index "github.com/blevesearch/bleve_index_api"
)

type VecPostingsList interface {
	DiskStatsReporter

	Iterator(prealloc VecPostingsIterator) VecPostingsIterator

	Size() int

	Count() uint64

	// NOTE deferred for future work

	// And(other PostingsList) PostingsList
	// Or(other PostingsList) PostingsList
}

type VecPostingsIterator interface {
	DiskStatsReporter

	// The caller is responsible for copying whatever it needs from
	// the returned Posting instance before calling Next(), as some
	// implementations may return a shared instance to reduce memory
	// allocations.
	Next() (VecPosting, error)

	// Advance will return the posting with the specified doc number
	// or if there is no such posting, the next posting.
	// Callers MUST NOT attempt to pass a docNum that is less than or
	// equal to the currently visited posting doc Num.
	Advance(docNum uint64) (VecPosting, error)

	Size() int
}

type VectorIndex interface {
	// Search performs a kNN search for the given query vector and returns a postings list.
	// - qVector: the query vector
	// - k: the number of similar vectors to return
	// - params: additional search parameters
	Search(qVector []float32, k int64, params json.RawMessage) (VecPostingsList, error)
	// SearchWithFilter performs a kNN search for the given query vector, filtering results based on eligible documents
	// - qVector: the query vector
	// - k: the number of similar vectors to return
	// - eligibleList: list of eligible documents to consider
	// - params: additional search parameters
	SearchWithFilter(qVector []float32, k int64, eligibleList index.EligibleDocumentList, params json.RawMessage) (VecPostingsList, error)
	// Close releases any resources held by the VectorIndex.
	Close()
	Size() uint64

	ObtainKCentroidCardinalitiesFromIVFIndex(limit int, descending bool) ([]index.CentroidCardinality, error)
}

type TrainedSegment interface {
	Segment
	GetCoarseQuantizer(field string) (interface{}, error)
}

type VectorSegment interface {
	Segment
	InterpretVectorIndex(field string, except *roaring.Bitmap) (VectorIndex, error)
}

type VecPosting interface {
	Number() uint64

	Score() float32

	Size() int
}
