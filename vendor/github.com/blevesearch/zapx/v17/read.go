//  Copyright (c) 2017 Couchbase, Inc.
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

import "encoding/binary"

func (sb *SegmentBase) getDocStoredMetaAndCompressed(docNum uint64) ([]byte, []byte, error) {
	_, storedOffset, n, metaLen, dataLen := sb.getDocStoredOffsets(docNum)

	meta := sb.mem[storedOffset+n : storedOffset+n+metaLen]
	data := sb.mem[storedOffset+n+metaLen : storedOffset+n+metaLen+dataLen]

	meta, err := sb.fileReader.process(meta)
	if err != nil {
		return nil, nil, err
	}
	data, err = sb.fileReader.process(data)
	if err != nil {
		return nil, nil, err
	}

	return meta, data, nil
}

func (sb *SegmentBase) getDocStoredOffsets(docNum uint64) (
	uint64, uint64, uint64, uint64, uint64) {
	indexOffset := sb.storedIndexOffset + (8 * docNum)

	storedOffset := binary.BigEndian.Uint64(sb.mem[indexOffset : indexOffset+8])

	var n uint64

	metaLen, read := binary.Uvarint(sb.mem[storedOffset : storedOffset+binary.MaxVarintLen64])
	n += uint64(read)

	dataLen, read := binary.Uvarint(sb.mem[storedOffset+n : storedOffset+n+binary.MaxVarintLen64])
	n += uint64(read)

	return indexOffset, storedOffset, n, metaLen, dataLen
}

func (sb *SegmentBase) getEdgeListOffset() uint64 {
	// if no stored index, then no edge list
	if sb.storedIndexOffset == 0 {
		return 0
	}
	// Edge list comes after the stored fields index (doc stored offsets)
	// The stored index offset points to where the doc offsets start
	// So edge list starts right after the last document offset
	// which is at sb.storedIndexOffset + (8 * sb.numDocs)
	// since each doc offset is 8 bytes
	return sb.storedIndexOffset + (8 * sb.numDocs)
}
