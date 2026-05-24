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
	"errors"
	"fmt"
	"math"
	"sync/atomic"

	"github.com/RoaringBitmap/roaring/v2"
	index "github.com/blevesearch/bleve_index_api"
	faiss "github.com/blevesearch/go-faiss"
	seg "github.com/blevesearch/scorch_segment_api/v2"
)

func init() {
	registerSegmentSection(SectionFaissVectorIndex, &faissVectorIndexSection{})
	invertedTextIndexSectionExclusionChecks = append(invertedTextIndexSectionExclusionChecks, func(field index.Field) bool {
		_, ok := field.(index.VectorField)
		return ok
	})
	faiss.SetOMPThreads(defaultFaissOMPThreads)
}

const (
	// Set the default number of OMP threads to be used by FAISS
	// to 1 since openMP does not support goroutine based threading well.
	defaultFaissOMPThreads = 1
	// Divide the estimated nprobe with this value to optimize
	// for latency.
	nprobeLatencyOptimization = 2
	// The threshold for number of vectors beyond which we start building the ivf class
	// of indexes
	ivfThreshold = 1000
	// The threshold for number of vectors beyond which we consider fast merging
	// using faiss's native merge capabilities, instead of reconstructing and adding
	// vectors one by one
	ivfSq8Threshold = 10000
)

// Vector index types supported.
type faissIndexType uint64

const (
	// faissFP32Index represents the standard float32 index in Faiss,
	// which stores vectors in 32-bit floating point format.
	faissFP32Index faissIndexType = iota
	// faissBIVFIndex represents the binary IVF index in Faiss,
	// which stores vectors in a 8-bit binary format.
	faissBIVFIndex
)

// Errors for invariant violations related to fast merge and trained index retrieval.
var (
	ErrorTrainedIndexNotIVF   error = errors.New("trained index is not an IVF index, which is required for fast merge")
	ErrorFastMergeIndexNotIVF error = errors.New("fast merge is only supported for IVF indexes")
)

type faissVectorIndexSection struct {
}

func (v *faissVectorIndexSection) Process(opaque map[int]resetable, docNum uint32, field index.Field, fieldID uint16) {
	if fieldID == math.MaxUint16 {
		return
	}
	if vf, ok := field.(index.VectorField); ok {
		vo := v.getVectorIndexOpaque(opaque)
		vo.process(vf, field.Name(), fieldID, docNum)
	}
}

func (v *faissVectorIndexSection) Persist(opaque map[int]resetable, w *FileWriter) error {
	vo := v.getVectorIndexOpaque(opaque)
	return vo.writeVectorIndexes(w)
}

func (v *faissVectorIndexSection) AddrForField(opaque map[int]resetable, fieldID int) int {
	vo := v.getVectorIndexOpaque(opaque)
	return vo.fieldAddrs[uint16(fieldID)]
}

// vecIndexInfo contains information specific to a vector index,
// including metadata and the faiss index pointer itself.
type vecIndexInfo struct {
	startOffset       int
	indexSize         uint64
	vecIds            []int64
	indexOptimizedFor string
	indexType         faissIndexType
	index             faissIndex
}

// Merge merges vector indexes from multiple segments into a single index.
func (v *faissVectorIndexSection) Merge(opaque map[int]resetable, segments []*SegmentBase,
	drops []*roaring.Bitmap, fieldsInv []string, newDocNumsIn [][]uint64, w *FileWriter,
	closeCh chan struct{}) error {
	vo := v.getVectorIndexOpaque(opaque)
	// preallocating the space over here, if there are too many fields
	// in the segment this will help by avoiding multiple allocation
	// calls.
	// the segments with valid vector sections in them
	vecSegs := make([]*SegmentBase, 0, len(segments))
	// vector index information from those segments
	indexes := make([]*vecIndexInfo, 0, len(segments))
	// mapping from vector IDs to docIDs across segments
	vecToDocID := make([]uint64, 0, len(segments))
	// for every field, gather the vector indexes from the segments
	// that have them, merge them and write them out to the writer.
	for fieldID, fieldName := range fieldsInv {
		// continue if field is not required to be indexed
		if !vo.fieldsOptions[fieldName].IsIndexed() {
			continue
		}
		indexes = indexes[:0] // resizing the slices
		vecSegs = vecSegs[:0]
		vecToDocID = vecToDocID[:0]
		// flag to indicate if there are deleted/updated vectors
		// in any of the vector indexes being merged.
		var drops bool
		for segI, sb := range segments {
			if isClosed(closeCh) {
				return seg.ErrClosed
			}
			if _, ok := sb.fieldsMap[fieldName]; !ok {
				continue
			}
			// check if the section address is a valid one for "fieldName" in the
			// segment sb. the local fieldID (fetched by the fieldsMap of the sb)
			// is to be used while consulting the fieldsSectionsMap
			pos := int(sb.fieldsSectionsMap[sb.fieldsMap[fieldName]-1][SectionFaissVectorIndex])
			if pos == 0 {
				continue
			}

			// loading doc values - adhering to the sections format. never
			// valid values for vector section
			_, n := binary.Uvarint(sb.mem[pos : pos+binary.MaxVarintLen64])
			pos += n
			_, n = binary.Uvarint(sb.mem[pos : pos+binary.MaxVarintLen64])
			pos += n

			// read the vector index optimization type represented as an int
			indexOptimizationTypeInt, n := binary.Uvarint(sb.mem[pos : pos+binary.MaxVarintLen64])
			pos += n
			// read the number of vectors
			numVecs, n := binary.Uvarint(sb.mem[pos : pos+binary.MaxVarintLen64])
			pos += n

			// track the valid vectors to be reconstructed for this segment
			// during the merge operation.
			newIndexInfo := &vecIndexInfo{
				indexOptimizedFor: index.VectorIndexOptimizationsReverseLookup[int(indexOptimizationTypeInt)],
				vecIds:            make([]int64, 0, numVecs),
			}

			// read the length of the docID list
			listLen, n := binary.Uvarint(sb.mem[pos : pos+binary.MaxVarintLen64])
			pos += n
			buf, err := sb.fileReader.process(sb.mem[pos : pos+int(listLen)])
			if err != nil {
				return err
			}
			pos += int(listLen)

			bufPos := 0
			bufLen := len(buf)
			for vecID := 0; vecID < int(numVecs); vecID++ {
				docID, n := binary.Uvarint(buf[bufPos:min(bufPos+binary.MaxVarintLen64, bufLen)])
				bufPos += n
				// check if this docID is dropped in the new segment
				newDocID := newDocNumsIn[segI][uint32(docID)]
				if newDocID != docDropped {
					// valid docID, track the mapping
					vecToDocID = append(vecToDocID, newDocID)
					// if the remapped doc ID is valid, track it
					// as part of vecs to be reconstructed (for larger indexes).
					// This accounts only for valid vector IDs, so deleted
					// ones won't be reconstructed in the final index.
					newIndexInfo.vecIds = append(newIndexInfo.vecIds, int64(vecID))
				} else {
					// some vectors are dropped, so we can't do a fast merge using faiss's
					// native merge capabilities, because of the data drift issue.
					drops = true
				}
			}

			if len(newIndexInfo.vecIds) == 0 {
				// no valid vectors to be merged from this segment
				continue
			}

			// read the type of vector index
			indexType, n := binary.Uvarint(sb.mem[pos : pos+binary.MaxVarintLen64])
			pos += n
			// read the size of the vector index
			indexSize, n := binary.Uvarint(sb.mem[pos : pos+binary.MaxVarintLen64])
			pos += n
			// record the start offset and size of the vector index
			newIndexInfo.startOffset = pos
			newIndexInfo.indexSize = indexSize
			newIndexInfo.indexType = faissIndexType(indexType)
			vecSegs = append(vecSegs, sb)
			indexes = append(indexes, newIndexInfo)
			pos += int(indexSize)
		}
		// continue if there are absolutely no valid vectors present in the segment
		// for this field and crucially don't store the section start offset in it
		if len(indexes) == 0 || len(vecToDocID) == 0 {
			continue
		}

		err := vo.flushSectionMetadata(fieldID, w, vecToDocID, indexes)
		if err != nil {
			return err
		}

		// we're going to use the trained index template regardless of whether there's
		// a update/delete in the segments being merged and we let the fast merge
		// path handle what exactly to do in such cases.
		trainedIndex, err := trainedIndexFromConfig(vo.config, fieldName)
		if err != nil {
			return err
		}

		useGPU := vo.fieldsOptions[fieldName].UseGPU()
		err = vo.mergeAndWriteVectorIndexes(trainedIndex, vecSegs, indexes, w, closeCh, useGPU, drops)
		if err != nil {
			return err
		}
	}
	return nil
}

func trainedIndexFromConfig(config map[string]interface{}, fieldName string) (faissIndexIVF, error) {
	var trainedIndexFor index.TrainedIndexCallbackFn
	var training bool
	var rv faissIndex
	if cb, ok := config[index.TrainedIndexCallback]; ok {
		trainedIndexFor = cb.(index.TrainedIndexCallbackFn)
	}
	if tf, ok := config[index.TrainingKey]; ok {
		training = tf.(bool)
	}
	// if we have a callback registered AND if the training flag is not set:
	//  - fastmerge is supported for this index
	// 	- we're not in the training phase of index creation where you want to be
	//    able to reconstruct the vectors for training
	if trainedIndexFor != nil && !training {
		trainedIndex, err := trainedIndexFor(fieldName)
		if err != nil {
			return nil, err
		}
		if trainedIndex != nil {
			rv = trainedIndex.(faissIndex)
		}
	}
	if rv == nil {
		return nil, nil
	}

	trainedIndex := rv.castIVF()
	if trainedIndex == nil {
		return nil, ErrorTrainedIndexNotIVF
	}
	return trainedIndex, nil
}

func (v *vectorIndexOpaque) flushSectionMetadata(fieldID int, w *FileWriter,
	vecToDocID []uint64, indexes []*vecIndexInfo) error {
	tempBuf := v.grabBuf(binary.MaxVarintLen64)
	fieldStart := w.Count()

	// marking the fact that for vector index, doc values are not valid by
	// storing fieldNotUninverted values.
	n := binary.PutUvarint(tempBuf, fieldNotUninverted)
	_, err := w.Write(tempBuf[:n])
	if err != nil {
		return err
	}
	n = binary.PutUvarint(tempBuf, fieldNotUninverted)
	_, err = w.Write(tempBuf[:n])
	if err != nil {
		return err
	}

	// write the index optimization type
	n = binary.PutUvarint(tempBuf, uint64(index.SupportedVectorIndexOptimizations[indexes[0].indexOptimizedFor]))
	_, err = w.Write(tempBuf[:n])
	if err != nil {
		return err
	}
	// write the number of vectors
	n = binary.PutUvarint(tempBuf, uint64(len(vecToDocID)))
	_, err = w.Write(tempBuf[:n])
	if err != nil {
		return err
	}

	buf := make([]byte, binary.MaxVarintLen64*len(vecToDocID))
	bufPos := 0
	for _, docID := range vecToDocID {
		n = binary.PutUvarint(buf[bufPos:], docID)
		bufPos += n
	}
	buf = w.process(buf[:bufPos])

	// write the size of the vector to docID map
	n = binary.PutUvarint(tempBuf, uint64(len(buf)))
	_, err = w.Write(tempBuf[:n])
	if err != nil {
		return err
	}
	// write the vecID -> docID mapping
	_, err = w.Write(buf)
	if err != nil {
		return err
	}

	// record the fieldStart value for this section.
	v.fieldAddrs[uint16(fieldID)] = fieldStart
	return nil
}

// Calculates the nprobe count, given nlist(number of centroids) based on
// the metric the index is optimized for.
func calculateNprobe(nlist int, indexOptimizedFor string) int32 {
	nprobe := int32(math.Sqrt(float64(nlist)))
	if indexOptimizedFor == index.IndexOptimizedForLatency {
		nprobe /= nprobeLatencyOptimization
		if nprobe < 1 {
			nprobe = 1
		}
	}
	return nprobe
}

// todo: need to detect and handle data drift in a more intelligent way
func (v *vectorIndexOpaque) fastMergeIndexes(trainedIndex faissIndexIVF, cfg *faissIndexConfig,
	drops bool, vecIndexes []*vecIndexInfo, w *FileWriter, closeCh chan struct{}) error {
	// create a faissIndex for merged index using nlist and nprobe from trained index's
	// config and we're hitting the fast merge path only if we've not enabled GPU
	nprobe, nlist := trainedIndex.ivfParams()
	cfg.nlist = nlist
	mergedIdx, err := faissIndexFactory(cfg)
	if err != nil {
		return err
	}
	defer mergedIdx.close()

	// cast to IVF index to be able to set the quantizer for the fast merge
	ivfMergedIdx := mergedIdx.castIVF()
	if ivfMergedIdx == nil {
		return ErrorFastMergeIndexNotIVF
	}
	err = ivfMergedIdx.setDirectMap(1)
	if err != nil {
		return err
	}
	// setting the same nprobe value in the merged index as the centroid
	// index to ensure that we probe the same number of clusters
	ivfMergedIdx.setNProbe(int32(nprobe))
	// using the trained index to copy the quantizers and set them in the final
	// merged index if possible
	err = ivfMergedIdx.setQuantizers(trainedIndex)
	if err != nil {
		return err
	}

	var reconsVecs []float32
	reconsAndAddFrom := func(vecIDs []int64, srcIdx faissIndex) error {
		neededReconsLen := len(vecIDs) * cfg.dimension
		if cap(reconsVecs) < neededReconsLen {
			reconsVecs = make([]float32, neededReconsLen)
		}
		reconsVecs = reconsVecs[:neededReconsLen]
		reconsVecs, err = srcIdx.reconstructBatch(vecIDs, reconsVecs)
		if err != nil {
			return err
		}

		vecSet, err := newVectorSet(cfg.dimension, reconsVecs)
		if err != nil {
			return err
		}
		if cfg.indexType == faissBIVFIndex {
			vecSet.binarize()
		}
		// add to target index the reconstructed vectors for the valid vector IDs from the source index.
		err = mergedIdx.add(vecSet)
		if err != nil {
			return err
		}
		reconsVecs = reconsVecs[:0]
		return nil
	}

	for _, vi := range vecIndexes {
		if isClosed(closeCh) {
			return seg.ErrClosed
		}
		childIdx := vi.index
		if drops {
			// if there are some deletes or updates in the segments being merged,
			// we can't say definitely which mutation can cause a data drift solely
			// by the vector count - as a fallback mechanism we reconstruct + add
			// the vectors in this scenario using the trained template since we can't
			// use merge_from.
			err = reconsAndAddFrom(vi.vecIds, childIdx)
			if err != nil {
				return err
			}
		} else {
			if err = ivfMergedIdx.mergeFrom(childIdx, mergedIdx.ntotal()); err != nil {
				// either the childIdx isn't compatible for fast merge or merge_from failed
				// so, in either case we can fallback to reconstructing and adding the vectors
				// one by one from the source index to the target index as a error handling mechanism.
				err = reconsAndAddFrom(vi.vecIds, childIdx)
				if err != nil {
					return err
				}
			}
		}
	}

	tempBuf := v.grabBuf(binary.MaxVarintLen64)
	// write the type of the vector index
	n := binary.PutUvarint(tempBuf, uint64(cfg.indexType))
	_, err = w.Write(tempBuf[:n])
	if err != nil {
		return err
	}

	return mergedIdx.write(tempBuf, w)
}

func (v *vectorIndexOpaque) mergeAndWriteVectorIndexes(trainedIndex faissIndexIVF, sbs []*SegmentBase,
	vecIndexes []*vecIndexInfo, w *FileWriter, closeCh chan struct{}, useGPU, drops bool) error {
	// safe to assume that all the indexes are of the same config values, given
	// that they are extracted from the field mapping info.
	var dims, metric, indexDataCap, reconsCap, nvecs int
	var indexOptimizedFor string
	var indexType faissIndexType
	var validMerge bool

	for segI, segBase := range sbs {
		// Considering merge operations on vector indexes are expensive, it is
		// worth including an early exit if the merge is aborted, saving us
		// the resource spikes, even if temporary.
		if isClosed(closeCh) {
			freeReconstructedIndexes(vecIndexes)
			return seg.ErrClosed
		}
		// track which index we are currently processing
		currVecIndex := vecIndexes[segI]
		currNumVecs := len(currVecIndex.vecIds)
		// if no valid vectors for this index, don't bring it into memory
		if currNumVecs == 0 {
			continue
		}

		// read the serialized index bytes
		indexBytes, err := segBase.fileReader.process(segBase.mem[currVecIndex.startOffset : currVecIndex.startOffset+int(currVecIndex.indexSize)])
		if err != nil {
			freeReconstructedIndexes(vecIndexes)
			return err
		}
		ioFlags := faissIOFlags
		if trainedIndex == nil {
			ioFlags = faissIOFlagsReadOnly
		}
		// reconstruct the faiss index from the bytes
		faissIndex, err := faiss.ReadIndexFromBuffer(indexBytes, ioFlags)
		if err != nil {
			freeReconstructedIndexes(vecIndexes)
			return err
		}

		// set the dims and metric values from the constructed index.
		dims = faissIndex.D()
		// at least one valid index to be merged, mark the merge as valid.
		validMerge = true
		metric = faissIndex.MetricType()
		indexOptimizedFor = currVecIndex.indexOptimizedFor
		indexType = currVecIndex.indexType
		// update trackers for buffer capacities
		indexReconsLen := currNumVecs * dims
		if indexReconsLen > reconsCap {
			reconsCap = indexReconsLen
		}
		indexDataCap += indexReconsLen

		// track the reconstruct index for this vector index, which will be used
		// to reconstruct the vectors corresponding to the valid vector IDs for this index.
		config := newFaissIndexConfig(indexType, indexOptimizedFor, dims, metric, currNumVecs, determineCentroids(currNumVecs), useGPU)
		fIndex, err := newFaissFloat32IndexWithConfig(faissIndex, config)
		if err != nil {
			freeReconstructedIndexes(vecIndexes)
			return err
		}
		vecIndexes[segI].index = fIndex

		// load binary index from disk if present
		if currVecIndex.indexType == faissBIVFIndex {
			// get to the bivf part of the vector index section
			pos := currVecIndex.startOffset + int(currVecIndex.indexSize)
			binSize, n := binary.Uvarint(segBase.mem[pos : pos+binary.MaxVarintLen64])
			pos += n
			indexBytes, err = segBase.fileReader.process(segBase.mem[pos : pos+int(binSize)])
			if err != nil {
				freeReconstructedIndexes(vecIndexes)
				return err
			}

			binaryIndex, err := faiss.ReadBinaryIndexFromBuffer(indexBytes, ioFlags)
			if err != nil {
				freeReconstructedIndexes(vecIndexes)
				return err
			}
			vecIndexes[segI].index, err = newFaissBinaryIndexWithConfig(binaryIndex, faissIndex, config)
			if err != nil {
				freeReconstructedIndexes(vecIndexes)
				return err
			}
		}
		nvecs += currNumVecs
	}

	// not a valid merge operation as there are no valid indexes to merge.
	if !validMerge {
		return nil
	}
	// if no valid vectors after merge, nothing to do
	if nvecs == 0 {
		// no valid vectors for this index, so we don't even have to
		// record it in the section
		freeReconstructedIndexes(vecIndexes)
		return nil
	}

	// create the faiss index to hold the merged data, either via fast merge or reconstruction
	config := newFaissIndexConfig(indexType, indexOptimizedFor, dims, metric, nvecs, determineCentroids(nvecs), useGPU)
	// we perform fast merge if we're not using the GPU and if the trained index
	// is compatible to be used for fast merge
	if !useGPU && canFastMerge(trainedIndex, indexOptimizedFor, nvecs) {
		err := v.fastMergeIndexes(trainedIndex, config, drops, vecIndexes, w, closeCh)
		if err != nil {
			return err
		}
		// free the indexes as we won't need them anymore after the fast merge
		freeReconstructedIndexes(vecIndexes)
		return nil
	}

	// Reconstruct Merge Path:
	// merging of indexes with reconstruction method.
	// the vecIds in each index contain only the valid vectors,
	// so we reconstruct only those.
	indexData := make([]float32, 0, indexDataCap)
	// reusable buffer for reconstruction
	recons := make([]float32, 0, reconsCap)
	for idx, currVecIndex := range vecIndexes {
		if isClosed(closeCh) {
			freeReconstructedIndexes(vecIndexes)
			return seg.ErrClosed
		}
		currNumVecs := len(currVecIndex.vecIds)
		// reconstruct the vectors only if present, it could be that
		// some of the indexes had all of their vectors updated/deleted.
		if currNumVecs > 0 && vecIndexes[idx] != nil {
			neededReconsLen := currNumVecs * config.dimension
			recons = recons[:neededReconsLen]
			var err error
			fIndex := vecIndexes[idx].index
			recons, err = fIndex.reconstructBatch(currVecIndex.vecIds, recons)
			if err != nil {
				freeReconstructedIndexes(vecIndexes)
				return err
			}
			indexData = append(indexData, recons...)
		}
	}

	// freeing the reconstructed indexes immediately - waiting till the end
	// to do the same is not needed because the following operations don't need
	// the reconstructed ones anymore and doing so will hold up memory which can
	// be detrimental while creating indexes during introduction.
	freeReconstructedIndexes(vecIndexes)

	vecSet, err := newVectorSet(config.dimension, indexData)
	if err != nil {
		return err
	}
	return v.writeFaissIndex(vecSet, config, w)
}

// constructs a faiss on the vectors according to the provided config and writes it out
// the given writer
func (v *vectorIndexOpaque) writeFaissIndex(vecs *vectorSet, config *faissIndexConfig, w *FileWriter) error {
	// create the faiss index based on the provided description string, and the metric type.
	index, err := faissIndexFactory(config)
	if err != nil {
		return err
	}
	// ensure the faiss index is closed after use
	defer index.close()

	// binarize the vectors for BIVF indexes
	if config.indexType == faissBIVFIndex {
		vecs.binarize()
	}
	// if we are using an IVF index, train and add first, then set the direct map
	// and nprobe. The order matters for GPU indexes: CloneToCPU (done inside
	// trainAndAdd) clears the direct map and nprobe, so they must be set after.
	if ivfIndex := index.castIVF(); ivfIndex != nil {
		// train the vector index and add the vectors to it. The training step
		// performs k-means clustering to partition the data space such that during
		// search time we probe only a subset of vectors (non-exhaustive search).
		err = ivfIndex.trainAndAdd(vecs, vecs)
		if err != nil {
			return err
		}
		// the direct map maintained in the IVF index is essential for the
		// reconstruction of vectors based on the sequential vector IDs in the
		// future merges use direct map type 1 -> array based direct map, since
		// we have sequential vector IDs starting from 0 to N-1.
		err = ivfIndex.setDirectMap(1)
		if err != nil {
			return err
		}
		// calculate nprobe using a heuristic.
		nprobe := calculateNprobe(config.nlist, config.optimizationType)
		ivfIndex.setNProbe(nprobe)
	} else {
		// add the vectors to the index using sequential vector IDs starting
		// from 0 to N-1
		err = index.add(vecs)
		if err != nil {
			return err
		}
	}

	// get a temporary buffer for writing out the index
	tempBuf := v.grabBuf(binary.MaxVarintLen64)
	// write the type of the vector index
	n := binary.PutUvarint(tempBuf, uint64(config.indexType))
	_, err = w.Write(tempBuf[:n])
	if err != nil {
		return err
	}

	// serialize the merged index into a byte slice, and write it out
	err = index.write(tempBuf, w)
	if err != nil {
		return err
	}

	return nil
}

// returns the index description string and index type constant for the binary
// index to be created based on the number of vectors and centroids.
func determineBinaryIndexToUse(nvecs, nlist int) string {
	switch {
	case nvecs >= ivfThreshold:
		return fmt.Sprintf("BIVF%d", nlist)
	default:
		return "BFlat"
	}
}

// returns the index type constant for the vector index to be created based on the
// index optimization type specified in the field mapping.
func determineIndexTypeFromOptimization(indexOptimizedFor string) faissIndexType {
	if index.OptimizationRequiresBinaryIndex(indexOptimizedFor) {
		return faissBIVFIndex
	}
	return faissFP32Index
}

// freeReconstructedIndexes closes all faiss indexes in the provided slice.
func freeReconstructedIndexes(vecIndexes []*vecIndexInfo) {
	for _, entry := range vecIndexes {
		if entry != nil && entry.index != nil {
			entry.index.close()
		}
	}
}

// grabBuf returns a reusable buffer of the given size, allocating a new one if needed.
func (v *vectorIndexOpaque) grabBuf(size int) []byte {
	buf := v.tmp0
	if cap(buf) < size {
		buf = make([]byte, size)
		v.tmp0 = buf
	}
	return buf[:size]
}

// determineCentroids determines the number of centroids to use for an IVF index.
func determineCentroids(nvecs int) int {
	var nlist int
	switch {
	case nvecs >= 200000:
		nlist = int(4 * math.Sqrt(float64(nvecs)))
	case nvecs >= ivfThreshold:
		// 100 points per cluster is a reasonable default, considering the default
		// minimum and maximum points per cluster is 39 and 256 respectively.
		// Since it's a recommendation to have a minimum of 10 clusters, 1000(100 * 10)
		// was chosen as the lower threshold.
		nlist = nvecs / 100
	}
	return nlist
}

// determineFloat32IndexToUse returns a description string for the float32
// index and quantizer type, and an index type constant.
func determineFloat32IndexToUse(nvecs, nlist int, optimizationType string) string {
	if nvecs < ivfThreshold {
		return "Flat"
	}
	switch optimizationType {
	case index.IndexBIVFWithBackingFlat:
		return "Flat"
	case index.IndexBIVFWithBackingSQ8:
		return "SQ8"
	case index.IndexOptimizedForMemoryEfficient:
		return fmt.Sprintf("IVF%d,SQ4", nlist)
	case index.IndexIVFRaBitQ:
		return fmt.Sprintf("IVF%d,RaBitQ", nlist)
	default:
		switch {
		case nvecs >= ivfSq8Threshold:
			return fmt.Sprintf("IVF%d,SQ8", nlist)
		default:
			return fmt.Sprintf("IVF%d,Flat", nlist)
		}
	}
}

func (vo *vectorIndexOpaque) writeVectorIndexes(w *FileWriter) error {
	// for every fieldID, contents to store over here are:
	//    1. the serialized representation of the dense vector index.
	//    2. its constituent metadata like:
	//        a. number of vectors
	//        b. dimension of vectors
	//        c. distance metric
	//        d. index optimization type
	//        e. vectorID -> docID mapping
	tempBuf := vo.grabBuf(binary.MaxVarintLen64)
	for fieldID, content := range vo.fieldVectorIndex {
		// number of vectors to be indexed for this field
		nvecs := len(content.vecDocIDs)
		// Set the faiss metric type (default is Euclidean Distance or l2_norm)
		metric := faiss.MetricL2
		if content.metric == index.InnerProduct || content.metric == index.CosineSimilarity {
			// use the same FAISS metric for inner product and cosine similarity
			metric = faiss.MetricInnerProduct
		}

		// create a vector set wrapping the vector data
		vecSet, err := newVectorSet(content.dimension, content.vectors)
		if err != nil {
			return err
		}

		// record the fieldStart value for this section.
		fieldStart := w.Count()
		// writing out two offset values to indicate that the current field's
		// vector section doesn't have valid doc value content within it.
		n := binary.PutUvarint(tempBuf, fieldNotUninverted)
		_, err = w.Write(tempBuf[:n])
		if err != nil {
			return err
		}
		n = binary.PutUvarint(tempBuf, fieldNotUninverted)
		_, err = w.Write(tempBuf[:n])
		if err != nil {
			return err
		}

		// write the index optimization type
		n = binary.PutUvarint(tempBuf, uint64(index.SupportedVectorIndexOptimizations[content.optimizedFor]))
		_, err = w.Write(tempBuf[:n])
		if err != nil {
			return err
		}
		// write the number of vectors
		n = binary.PutUvarint(tempBuf, uint64(nvecs))
		_, err = w.Write(tempBuf[:n])
		if err != nil {
			return err
		}

		buf := make([]byte, binary.MaxVarintLen64*len(content.vecDocIDs))
		bufPos := 0
		for _, docID := range content.vecDocIDs {
			n = binary.PutUvarint(buf[bufPos:], uint64(docID))
			bufPos += n
		}
		buf = w.process(buf[:bufPos])

		// write the size of the vector to docID map
		n = binary.PutUvarint(tempBuf, uint64(len(buf)))
		_, err = w.Write(tempBuf[:n])
		if err != nil {
			return err
		}
		// write the vecID -> docID mapping
		_, err = w.Write(buf)
		if err != nil {
			return err
		}

		// determine the type of vector index to be created based on the index optimization
		// and create the faiss index for the vectors associated with this field and
		// write out the index into the segment writer.
		indexType := determineIndexTypeFromOptimization(content.optimizedFor)
		config := newFaissIndexConfig(indexType, content.optimizedFor, content.dimension, metric, nvecs, determineCentroids(nvecs), false)
		err = vo.writeFaissIndex(vecSet, config, w)
		if err != nil {
			return err
		}

		// accounts for whatever data has been written out to the writer.
		vo.incrementBytesWritten(uint64(w.Count() - fieldStart))
		vo.fieldAddrs[fieldID] = fieldStart
	}
	return nil
}

func (vo *vectorIndexOpaque) process(field index.VectorField, fieldName string, fieldID uint16, docNum uint32) {
	if fieldID == math.MaxUint16 {
		// doc processing checkpoint - no action needed
		return
	}
	vec := field.Vector()
	dim := field.Dims()
	metric := field.Similarity()
	indexOptimizedFor := field.IndexOptimizedFor()
	// caller is supposed to make sure len(vec) is a multiple of dim.
	// Not double checking it here to avoid the overhead.
	// This accounts for multi-vector fields, where a field can have
	// multiple vectors associated with it. In this case we process all
	// vectors associated with the field as separate vectors.
	numVectors := len(vec) / dim
	for i := 0; i < numVectors; i++ {
		vector := vec[i*dim : (i+1)*dim]
		// check if we have content for this fieldID already
		content, ok := vo.fieldVectorIndex[fieldID]
		if !ok {
			// create an entry for this fieldID as this is the first time
			// we are seeing this field
			content = &vectorIndexContent{
				dimension:    dim,
				metric:       metric,
				optimizedFor: indexOptimizedFor,
				vectors:      make([]float32, 0, dim*numVectors),
				vecDocIDs:    make([]uint32, 0, numVectors),
				useGPU:       vo.fieldsOptions[fieldName].UseGPU(),
			}
			vo.fieldVectorIndex[fieldID] = content
		}
		// track the vector data and docIDs
		content.vectors = append(content.vectors, vector...)
		content.vecDocIDs = append(content.vecDocIDs, docNum)
	}
}

func (v *faissVectorIndexSection) getVectorIndexOpaque(opaque map[int]resetable) *vectorIndexOpaque {
	if _, ok := opaque[SectionFaissVectorIndex]; !ok {
		opaque[SectionFaissVectorIndex] = v.InitOpaque(nil)
	}
	return opaque[SectionFaissVectorIndex].(*vectorIndexOpaque)
}

func (v *faissVectorIndexSection) InitOpaque(args map[string]interface{}) resetable {
	rv := &vectorIndexOpaque{
		fieldAddrs:       make(map[uint16]int),
		fieldVectorIndex: make(map[uint16]*vectorIndexContent),
	}
	for k, v := range args {
		rv.Set(k, v)
	}

	return rv
}

// vectorIndexContent contains the information required to create a vector index for a vector field.
type vectorIndexContent struct {
	// vectors stores flattened vectors in a row-major order
	vectors []float32
	// vecDocIDs corresponding to each vector
	vecDocIDs []uint32
	// dimension is the dimension of all vectors
	dimension int
	// metric is the distance metric to be used
	metric string
	// optimizedFor is the optimization type for the index
	optimizedFor string
	// useGPU indicates whether the index should be created on the GPU
	useGPU bool
}

// vectorIndexOpaque holds the internal state for vector index processing.
type vectorIndexOpaque struct {
	// external config values passed in, which controls the behavior of vector index creation and merging
	config map[string]interface{}
	// number of bytes written out for the vector index section, used for metrics and tracking
	bytesWritten uint64
	// fieldAddrs maps fieldID to the address of its vector section
	fieldAddrs map[uint16]int
	// fieldVectorIndex maps fieldID to its vector index content
	fieldVectorIndex map[uint16]*vectorIndexContent
	// fieldsOptions contains field indexing options
	fieldsOptions map[string]index.FieldIndexingOptions
	// tmp0 is a reusable buffer
	tmp0 []byte
}

func (vo *vectorIndexOpaque) incrementBytesWritten(val uint64) {
	atomic.AddUint64(&vo.bytesWritten, val)
}

func (vo *vectorIndexOpaque) BytesWritten() uint64 {
	return atomic.LoadUint64(&vo.bytesWritten)
}

func (vo *vectorIndexOpaque) BytesRead() uint64 {
	return 0
}

func (vo *vectorIndexOpaque) ResetBytesRead(uint64) {
}

// Reset clears all state in the vectorIndexOpaque for reuse.
func (vo *vectorIndexOpaque) Reset() error {
	clear(vo.fieldAddrs)
	clear(vo.fieldVectorIndex)
	vo.tmp0 = vo.tmp0[:0]
	vo.fieldsOptions = nil
	vo.config = nil
	atomic.StoreUint64(&vo.bytesWritten, 0)
	return nil
}

func (v *vectorIndexOpaque) Set(key string, val interface{}) {
	switch key {
	case "fieldsOptions":
		v.fieldsOptions = val.(map[string]index.FieldIndexingOptions)
	case "config":
		v.config = val.(map[string]interface{})
	}
}

// ---------------------------------
// Faiss Index Factory
// ---------------------------------
type faissIndexConfig struct {
	indexType        faissIndexType
	dimension        int
	metricType       int
	numVecs          int
	optimizationType string
	nlist            int
	useGPU           bool
}

func newFaissIndexConfig(idxType faissIndexType, optimizationType string, dimension, metricType, numVecs, nlist int, useGPU bool) *faissIndexConfig {
	return &faissIndexConfig{
		indexType:        idxType,
		dimension:        dimension,
		metricType:       metricType,
		numVecs:          numVecs,
		nlist:            nlist,
		optimizationType: optimizationType,
		useGPU:           useGPU,
	}
}

// Factory function to create a faissIndex for the given index config.
func faissIndexFactory(cfg *faissIndexConfig) (faissIndex, error) {
	switch cfg.indexType {
	case faissFP32Index:
		description := determineFloat32IndexToUse(cfg.numVecs, cfg.nlist, cfg.optimizationType)
		idx, err := faiss.IndexFactory(cfg.dimension, description, cfg.metricType)
		if err != nil {
			return nil, err
		}
		// we restrict GPU to IVF indexes only; flat and SQ indexes do not get a noticeable speedup
		// when run on GPU, and the GPU overhead can actually make them slower than CPU.
		if cfg.useGPU && idx.IsIVFIndex() {
			return newFaissGPUFloat32Index(idx)
		}
		return newFaissFloat32Index(idx)
	case faissBIVFIndex:
		description := determineBinaryIndexToUse(cfg.numVecs, cfg.nlist)
		binaryIdx, err := faiss.BinaryIndexFactory(cfg.dimension, description)
		if err != nil {
			return nil, err
		}

		description = determineFloat32IndexToUse(cfg.numVecs, cfg.nlist, cfg.optimizationType)
		backingIdx, err := faiss.IndexFactory(cfg.dimension, description, cfg.metricType)
		if err != nil {
			return nil, err
		}
		return newFaissBinaryIndex(binaryIdx, backingIdx)
	default:
		return nil, errNotSupported
	}
}

// canFastMerge determines whether we can use the fast merge capabilities of faiss based on
//   - the presence of a trained index
//   - the optimization type of the index.
//   - the total number of vectors being merged.
func canFastMerge(trainedIndex faissIndexIVF, opt string, totalVecs int) bool {
	// if the trained index isn't IVF or not available, fallback to naive merge
	if trainedIndex == nil {
		return false
	}

	var minVecsForFastMerge int
	switch opt {
	case index.IndexBIVFWithBackingFlat, index.IndexBIVFWithBackingSQ8:
		fallthrough
	case index.IndexOptimizedForMemoryEfficient:
		fallthrough
	case index.IndexIVFRaBitQ:
		minVecsForFastMerge = ivfThreshold
	default:
		minVecsForFastMerge = ivfSq8Threshold
	}
	return trainedIndex.ntotal() > int64(minVecsForFastMerge) && totalVecs > minVecsForFastMerge
}
