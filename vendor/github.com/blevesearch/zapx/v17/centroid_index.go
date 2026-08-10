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

//go:build vectors
// +build vectors

package zap

import (
	"encoding/binary"
	"fmt"

	faiss "github.com/blevesearch/go-faiss"
)

func (sb *SegmentBase) GetCoarseQuantizer(field string) (interface{}, error) {
	fieldIDPlus1 := sb.fieldsMap[field]
	if fieldIDPlus1 <= 0 {
		return nil, fmt.Errorf("field %s does not exist in segment", field)
	}

	vectorSection := sb.fieldsSectionsMap[fieldIDPlus1-1][SectionFaissVectorIndex]
	// check if the field has a vector section in the segment.
	if vectorSection <= 0 {
		return nil, fmt.Errorf("field %s does not have a vector section in the segment", field)
	}

	pos := int(vectorSection)
	// doc values and vector optimization type
	for i := 0; i < 3; i++ {
		_, n := binary.Uvarint(sb.mem[pos : pos+binary.MaxVarintLen64])
		pos += n
	}

	numVecs, n := binary.Uvarint(sb.mem[pos : pos+binary.MaxVarintLen64])
	pos += n

	// length of the vector to docID map
	_, n = binary.Uvarint(sb.mem[pos : pos+binary.MaxVarintLen64])
	pos += n

	// vector to docID mapping
	for i := 0; i < int(numVecs); i++ {
		_, n = binary.Uvarint(sb.mem[pos : pos+binary.MaxVarintLen64])
		pos += n
	}

	// type of index
	indexType, n := binary.Uvarint(sb.mem[pos : pos+binary.MaxVarintLen64])
	pos += n
	indexSize, n := binary.Uvarint(sb.mem[pos : pos+binary.MaxVarintLen64])
	pos += n

	// todo: might wanna use the vector cache here, early tests didn't show a big diff
	faissIndex, err := faiss.ReadIndexFromBuffer(sb.mem[pos:pos+int(indexSize)], faissIOFlags)
	if err != nil {
		return nil, err
	}
	pos += int(indexSize)

	if faissIndexType(indexType) == faissBIVFIndex {
		binaryIndexSize, n := binary.Uvarint(sb.mem[pos : pos+binary.MaxVarintLen64])
		pos += n
		binaryIndex, err := faiss.ReadBinaryIndexFromBuffer(sb.mem[pos:pos+int(binaryIndexSize)], faissIOFlags)
		if err != nil {
			return nil, err
		}
		return newFaissBinaryIndex(binaryIndex, faissIndex)
	}
	return newFaissFloat32Index(faissIndex)
}
