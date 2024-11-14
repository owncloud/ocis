//  Copyright (c) 2023 Couchbase, Inc.
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

//go:build vectors
// +build vectors

package zap

import (
	"encoding/binary"
	"encoding/json"
	"math"
	"reflect"

	"github.com/RoaringBitmap/roaring"
	"github.com/RoaringBitmap/roaring/roaring64"
	faiss "github.com/blevesearch/go-faiss"
	segment "github.com/blevesearch/scorch_segment_api/v2"
)

var reflectStaticSizeVecPostingsList int
var reflectStaticSizeVecPostingsIterator int
var reflectStaticSizeVecPosting int

func init() {
	var pl VecPostingsList
	reflectStaticSizeVecPostingsList = int(reflect.TypeOf(pl).Size())
	var pi VecPostingsIterator
	reflectStaticSizeVecPostingsIterator = int(reflect.TypeOf(pi).Size())
	var p VecPosting
	reflectStaticSizeVecPosting = int(reflect.TypeOf(p).Size())
}

type VecPosting struct {
	docNum uint64
	score  float32
}

func (vp *VecPosting) Number() uint64 {
	return vp.docNum
}

func (vp *VecPosting) Score() float32 {
	return vp.score
}

func (vp *VecPosting) Size() int {
	sizeInBytes := reflectStaticSizePosting

	return sizeInBytes
}

// =============================================================================

// the vector postings list is supposed to store the docNum and its similarity
// score as a vector postings entry in it.
// The way in which is it stored is using a roaring64 bitmap.
// the docNum is stored in high 32 and the lower 32 bits contains the score value.
// the score is actually a float32 value and in order to store it as a uint32 in
// the bitmap, we use the IEEE 754 floating point format.
//
// each entry in the roaring64 bitmap of the vector postings list is a 64 bit
// number which looks like this:
// MSB                         LSB
// |64 63 62 ... 32| 31 30 ... 0|
// |    <docNum>   |   <score>  |
type VecPostingsList struct {
	// todo: perhaps we don't even need to store a bitmap if there is only
	// one similar vector the query, but rather store it as a field value
	// in the struct
	except   *roaring64.Bitmap
	postings *roaring64.Bitmap
}

var emptyVecPostingsIterator = &VecPostingsIterator{}
var emptyVecPostingsList = &VecPostingsList{}

func (vpl *VecPostingsList) Iterator(prealloc segment.VecPostingsIterator) segment.VecPostingsIterator {

	// tbd: do we check the cardinality of postings and scores?
	var preallocPI *VecPostingsIterator
	pi, ok := prealloc.(*VecPostingsIterator)
	if ok && pi != nil {
		preallocPI = pi
	}
	if preallocPI == emptyVecPostingsIterator {
		preallocPI = nil
	}

	return vpl.iterator(preallocPI)
}

func (p *VecPostingsList) iterator(rv *VecPostingsIterator) *VecPostingsIterator {
	if rv == nil {
		rv = &VecPostingsIterator{}
	} else {
		*rv = VecPostingsIterator{} // clear the struct
	}
	// think on some of the edge cases over here.
	if p.postings == nil {
		return rv
	}
	rv.postings = p
	rv.all = p.postings.Iterator()
	if p.except != nil {
		rv.ActualBM = roaring64.AndNot(p.postings, p.except)
		rv.Actual = rv.ActualBM.Iterator()
	} else {
		rv.ActualBM = p.postings
		rv.Actual = rv.all // Optimize to use same iterator for all & Actual.
	}
	return rv
}

func (p *VecPostingsList) Size() int {
	sizeInBytes := reflectStaticSizeVecPostingsList + SizeOfPtr

	if p.except != nil {
		sizeInBytes += int(p.except.GetSizeInBytes())
	}

	return sizeInBytes
}

func (p *VecPostingsList) Count() uint64 {
	if p.postings != nil {
		n := p.postings.GetCardinality()
		var e uint64
		if p.except != nil {
			e = p.postings.AndCardinality(p.except)
		}
		return n - e
	}
	return 0
}

func (vpl *VecPostingsList) ResetBytesRead(val uint64) {

}

func (vpl *VecPostingsList) BytesRead() uint64 {
	return 0
}

func (vpl *VecPostingsList) BytesWritten() uint64 {
	return 0
}

// =============================================================================

type VecPostingsIterator struct {
	postings *VecPostingsList
	all      roaring64.IntPeekable64
	Actual   roaring64.IntPeekable64
	ActualBM *roaring64.Bitmap

	next VecPosting // reused across Next() calls
}

func (i *VecPostingsIterator) nextCodeAtOrAfterClean(atOrAfter uint64) (uint64, bool, error) {
	i.Actual.AdvanceIfNeeded(atOrAfter)

	if !i.Actual.HasNext() {
		return 0, false, nil // couldn't find anything
	}

	return i.Actual.Next(), true, nil
}

func (i *VecPostingsIterator) nextCodeAtOrAfter(atOrAfter uint64) (uint64, bool, error) {
	if i.Actual == nil || !i.Actual.HasNext() {
		return 0, false, nil
	}

	if i.postings == nil || i.postings == emptyVecPostingsList {
		// couldn't find anything
		return 0, false, nil
	}

	if i.postings.postings == i.ActualBM {
		return i.nextCodeAtOrAfterClean(atOrAfter)
	}

	i.Actual.AdvanceIfNeeded(atOrAfter)

	if !i.Actual.HasNext() || !i.all.HasNext() {
		// couldn't find anything
		return 0, false, nil
	}

	n := i.Actual.Next()
	allN := i.all.Next()

	// n is the next actual hit (excluding some postings), and
	// allN is the next hit in the full postings, and
	// if they don't match, move 'all' forwards until they do.
	for allN != n {
		if !i.all.HasNext() {
			return 0, false, nil
		}
		allN = i.all.Next()
	}

	return uint64(n), true, nil
}

// a transformation function which stores both the score and the docNum as a single
// entry which is a uint64 number.
func getVectorCode(docNum uint32, score float32) uint64 {
	return uint64(docNum)<<32 | uint64(math.Float32bits(score))
}

// Next returns the next posting on the vector postings list, or nil at the end
func (i *VecPostingsIterator) nextAtOrAfter(atOrAfter uint64) (segment.VecPosting, error) {
	// transform the docNum provided to the vector code format and use that to
	// get the next entry. the comparison still happens docNum wise since after
	// the transformation, the docNum occupies the upper 32 bits just an entry in
	// the postings list
	atOrAfter = getVectorCode(uint32(atOrAfter), 0)
	code, exists, err := i.nextCodeAtOrAfter(atOrAfter)
	if err != nil || !exists {
		return nil, err
	}

	i.next = VecPosting{} // clear the struct
	rv := &i.next
	rv.score = math.Float32frombits(uint32(code))
	rv.docNum = code >> 32

	return rv, nil
}

func (itr *VecPostingsIterator) Next() (segment.VecPosting, error) {
	return itr.nextAtOrAfter(0)
}

func (itr *VecPostingsIterator) Advance(docNum uint64) (segment.VecPosting, error) {
	return itr.nextAtOrAfter(docNum)
}

func (i *VecPostingsIterator) Size() int {
	sizeInBytes := reflectStaticSizePostingsIterator + SizeOfPtr +
		i.next.Size()

	return sizeInBytes
}

func (vpl *VecPostingsIterator) ResetBytesRead(val uint64) {

}

func (vpl *VecPostingsIterator) BytesRead() uint64 {
	return 0
}

func (vpl *VecPostingsIterator) BytesWritten() uint64 {
	return 0
}

// vectorIndexWrapper conforms to scorch_segment_api's VectorIndex interface
type vectorIndexWrapper struct {
	search func(qVector []float32, k int64,
		params json.RawMessage) (segment.VecPostingsList, error)
	searchWithFilter func(qVector []float32, k int64, eligibleDocIDs []uint64,
		params json.RawMessage) (segment.VecPostingsList, error)
	close func()
	size  func() uint64
}

func (i *vectorIndexWrapper) Search(qVector []float32, k int64,
	params json.RawMessage) (
	segment.VecPostingsList, error) {
	return i.search(qVector, k, params)
}

func (i *vectorIndexWrapper) SearchWithFilter(qVector []float32, k int64,
	eligibleDocIDs []uint64, params json.RawMessage) (
	segment.VecPostingsList, error) {
	return i.searchWithFilter(qVector, k, eligibleDocIDs, params)
}

func (i *vectorIndexWrapper) Close() {
	i.close()
}

func (i *vectorIndexWrapper) Size() uint64 {
	return i.size()
}

// InterpretVectorIndex returns a construct of closures (vectorIndexWrapper)
// that will allow the caller to -
// (1) search within an attached vector index
// (2) search limited to a subset of documents within an attached vector index
// (3) close attached vector index
// (4) get the size of the attached vector index
func (sb *SegmentBase) InterpretVectorIndex(field string, requiresFiltering bool,
	except *roaring.Bitmap) (
	segment.VectorIndex, error) {
	// Params needed for the closures
	var vecIndex *faiss.IndexImpl
	var vecDocIDMap map[int64]uint32
	var docVecIDMap map[uint32][]int64
	var vectorIDsToExclude []int64
	var fieldIDPlus1 uint16
	var vecIndexSize uint64

	// Utility function to add the corresponding docID and scores for each vector
	// returned after the kNN query to the newly
	// created vecPostingsList
	addIDsToPostingsList := func(pl *VecPostingsList, ids []int64, scores []float32) {
		for i := 0; i < len(ids); i++ {
			vecID := ids[i]
			// Checking if it's present in the vecDocIDMap.
			// If -1 is returned as an ID(insufficient vectors), this will ensure
			// it isn't added to the final postings list.
			if docID, ok := vecDocIDMap[vecID]; ok {
				code := getVectorCode(docID, scores[i])
				pl.postings.Add(uint64(code))
			}
		}
	}

	var (
		wrapVecIndex = &vectorIndexWrapper{
			search: func(qVector []float32, k int64, params json.RawMessage) (
				segment.VecPostingsList, error) {
				// 1. returned postings list (of type PostingsList) has two types of information - docNum and its score.
				// 2. both the values can be represented using roaring bitmaps.
				// 3. the Iterator (of type PostingsIterator) returned would operate in terms of VecPostings.
				// 4. VecPostings would just have the docNum and the score. Every call of Next()
				//    and Advance just returns the next VecPostings. The caller would do a vp.Number()
				//    and the Score() to get the corresponding values
				rv := &VecPostingsList{
					except:   nil, // todo: handle the except bitmap within postings iterator.
					postings: roaring64.New(),
				}

				if vecIndex == nil || vecIndex.D() != len(qVector) {
					// vector index not found or dimensionality mismatched
					return rv, nil
				}

				scores, ids, err := vecIndex.SearchWithoutIDs(qVector, k,
					vectorIDsToExclude, params)
				if err != nil {
					return nil, err
				}

				addIDsToPostingsList(rv, ids, scores)

				return rv, nil
			},
			searchWithFilter: func(qVector []float32, k int64,
				eligibleDocIDs []uint64, params json.RawMessage) (
				segment.VecPostingsList, error) {
				// 1. returned postings list (of type PostingsList) has two types of information - docNum and its score.
				// 2. both the values can be represented using roaring bitmaps.
				// 3. the Iterator (of type PostingsIterator) returned would operate in terms of VecPostings.
				// 4. VecPostings would just have the docNum and the score. Every call of Next()
				//    and Advance just returns the next VecPostings. The caller would do a vp.Number()
				//    and the Score() to get the corresponding values
				rv := &VecPostingsList{
					except:   nil, // todo: handle the except bitmap within postings iterator.
					postings: roaring64.New(),
				}

				if vecIndex == nil || vecIndex.D() != len(qVector) {
					// vector index not found or dimensionality mismatched
					return rv, nil
				}

				if len(eligibleDocIDs) > 0 {
					// Non-zero documents eligible per the filter query.

					// If every element in the index is eligible(eg. high selectivity
					// cases), then this can basically be considered unfiltered kNN.
					if len(eligibleDocIDs) == int(sb.numDocs) {
						scores, ids, err := vecIndex.SearchWithoutIDs(qVector, k,
							vectorIDsToExclude, params)
						if err != nil {
							return nil, err
						}

						addIDsToPostingsList(rv, ids, scores)
						return rv, nil
					}

					// vector IDs corresponding to the local doc numbers to be
					// considered for the search
					vectorIDsToInclude := make([]int64, 0, len(eligibleDocIDs))
					for _, id := range eligibleDocIDs {
						vectorIDsToInclude = append(vectorIDsToInclude, docVecIDMap[uint32(id)]...)
					}

					// Retrieve the mapping of centroid IDs to vectors within
					// the cluster.
					clusterAssignment, _ := vecIndex.ObtainClusterToVecIDsFromIVFIndex()
					// Accounting for a flat index
					if len(clusterAssignment) == 0 {
						scores, ids, err := vecIndex.SearchWithIDs(qVector, k,
							vectorIDsToInclude, params)
						if err != nil {
							return nil, err
						}

						addIDsToPostingsList(rv, ids, scores)
						return rv, nil
					}

					// Converting to roaring bitmap for ease of intersect ops with
					// the set of eligible doc IDs.
					centroidVecIDMap := make(map[int64]*roaring.Bitmap)
					for centroidID, vecIDs := range clusterAssignment {
						if _, exists := centroidVecIDMap[centroidID]; !exists {
							centroidVecIDMap[centroidID] = roaring.NewBitmap()
						}
						vecIDsUint32 := make([]uint32, 0, len(vecIDs))
						for _, vecID := range vecIDs {
							vecIDsUint32 = append(vecIDsUint32, uint32(vecID))
						}
						centroidVecIDMap[centroidID].AddMany(vecIDsUint32)
					}

					// Determining which clusters, identified by centroid ID,
					// have at least one eligible vector and hence, ought to be
					// probed.
					eligibleCentroidIDs := make([]int64, 0)

					var selector faiss.Selector
					var err error
					// If there are more elements to be included than excluded, it
					// might be quicker to use an exclusion selector as a filter
					// instead of an inclusion selector.
					if float32(len(eligibleDocIDs))/float32(len(docVecIDMap)) > 0.5 {
						ineligibleVecIDsBitmap := roaring.NewBitmap()
						eligibleDocIDsMap := make(map[uint64]struct{})
						for _, eligibleDocID := range eligibleDocIDs {
							eligibleDocIDsMap[(eligibleDocID)] = struct{}{}
						}

						ineligibleVectorIDs := make([]int64, 0, len(vecDocIDMap)-
							len(vectorIDsToInclude))

						for docID, vecIDs := range docVecIDMap {
							if _, exists := eligibleDocIDsMap[uint64(docID)]; !exists {
								for _, vecID := range vecIDs {
									ineligibleVecIDsBitmap.Add(uint32(vecID))
									ineligibleVectorIDs = append(ineligibleVectorIDs, vecID)
								}
							}
						}

						for centroidID, vecIDs := range centroidVecIDMap {
							vecIDs.AndNot(ineligibleVecIDsBitmap)
							// At least one eligible vec in cluster.
							if !vecIDs.IsEmpty() {
								// The mapping is now reduced to those vectors which
								// are also eligible docs for the filter query.
								centroidVecIDMap[centroidID] = vecIDs
								eligibleCentroidIDs = append(eligibleCentroidIDs, centroidID)
							} else {
								// don't consider clusters with no eligible IDs.
								delete(centroidVecIDMap, centroidID)
							}
						}

						selector, err = faiss.NewIDSelectorNot(ineligibleVectorIDs)
					} else {
						// Getting the vector IDs corresponding to the eligible
						// doc IDs.
						// The docVecIDMap maps each docID to vectorIDs corresponding
						// to it.
						// Usually, each docID has one vecID mapped to it unless
						// the vector is nested, in which case there can be multiple
						// vectorIDs mapped to the same docID.
						// Eg. docID d1 -> vecID v1, for the first case
						// d1 -> {v1,v2}, for the second case.
						eligibleVecIDsBitmap := roaring.NewBitmap()
						vecIDsUint32 := make([]uint32, 0)
						for _, eligibleDocID := range eligibleDocIDs {
							vecIDs := docVecIDMap[uint32(eligibleDocID)]
							for _, vecID := range vecIDs {
								vecIDsUint32 = append(vecIDsUint32, uint32(vecID))
							}
						}
						eligibleVecIDsBitmap.AddMany(vecIDsUint32)
						for centroidID, vecIDs := range centroidVecIDMap {
							vecIDs.And(eligibleVecIDsBitmap)
							if !vecIDs.IsEmpty() {
								// The mapping is now reduced to those vectors which
								// are also eligible docs for the filter query.
								centroidVecIDMap[centroidID] = vecIDs
								eligibleCentroidIDs = append(eligibleCentroidIDs, centroidID)
							} else {
								// don't consider clusters with no eligible IDs.
								delete(centroidVecIDMap, centroidID)
							}
						}

						selector, err = faiss.NewIDSelectorBatch(vectorIDsToInclude)
					}
					if err != nil {
						return nil, err
					}

					// Ordering the retrieved centroid IDs by increasing order
					// of distance i.e. decreasing order of proximity to query vector.
					closestCentroidIDs, centroidDistances, _ :=
						vecIndex.ObtainClustersWithDistancesFromIVFIndex(qVector,
							eligibleCentroidIDs)

					// Getting the nprobe value set at index time.
					nprobe := vecIndex.GetNProbe()

					eligibleDocsTillNow := int64(0)
					minEligibleCentroids := 0
					for i, centroidID := range closestCentroidIDs {
						eligibleDocsTillNow += int64(centroidVecIDMap[centroidID].GetCardinality())
						if eligibleDocsTillNow >= k && i >= int(nprobe-1) {
							// Continue till at least 'K' cumulative vectors are
							// collected or 'nprobe' clusters are examined, whichever
							// comes later.
							minEligibleCentroids = i + 1
							break
						}
						minEligibleCentroids = i + 1
					}

					// Search the clusters specified by 'closestCentroidIDs' for
					// vectors whose IDs are present in 'vectorIDsToInclude'
					scores, ids, err := vecIndex.SearchClustersFromIVFIndex(
						selector, len(vectorIDsToInclude), closestCentroidIDs,
						minEligibleCentroids, k, qVector, centroidDistances, params)
					if err != nil {
						return nil, err
					}

					addIDsToPostingsList(rv, ids, scores)
					return rv, nil
				}
				return rv, nil
			},
			close: func() {
				// skipping the closing because the index is cached and it's being
				// deferred to a later point of time.
				sb.vecIndexCache.decRef(fieldIDPlus1)
			},
			size: func() uint64 {
				return vecIndexSize
			},
		}

		err error
	)

	fieldIDPlus1 = sb.fieldsMap[field]
	if fieldIDPlus1 <= 0 {
		return wrapVecIndex, nil
	}

	vectorSection := sb.fieldsSectionsMap[fieldIDPlus1-1][SectionFaissVectorIndex]
	// check if the field has a vector section in the segment.
	if vectorSection <= 0 {
		return wrapVecIndex, nil
	}

	pos := int(vectorSection)

	// the below loop loads the following:
	// 1. doc values(first 2 iterations) - adhering to the sections format. never
	// valid values for vector section
	// 2. index optimization type.
	for i := 0; i < 3; i++ {
		_, n := binary.Uvarint(sb.mem[pos : pos+binary.MaxVarintLen64])
		pos += n
	}

	vecIndex, vecDocIDMap, docVecIDMap, vectorIDsToExclude, err =
		sb.vecIndexCache.loadOrCreate(fieldIDPlus1, sb.mem[pos:], requiresFiltering,
			except)

	if vecIndex != nil {
		vecIndexSize = vecIndex.Size()
	}

	return wrapVecIndex, err
}

func (sb *SegmentBase) UpdateFieldStats(stats segment.FieldStats) {
	for _, fieldName := range sb.fieldsInv {
		pos := int(sb.fieldsSectionsMap[sb.fieldsMap[fieldName]-1][SectionFaissVectorIndex])
		if pos == 0 {
			continue
		}

		for i := 0; i < 3; i++ {
			_, n := binary.Uvarint(sb.mem[pos : pos+binary.MaxVarintLen64])
			pos += n
		}
		numVecs, _ := binary.Uvarint(sb.mem[pos : pos+binary.MaxVarintLen64])

		stats.Store("num_vectors", fieldName, numVecs)
	}
}
