/*
 * Copyright 2019 Kopano and its licensors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *	http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package kcc

import (
	"bytes"
	"encoding/binary"
)

var (
	// MUIDECSAB is the GUID used in AB EntryIDs (ABEID). Definition copied
	// from kopanocore/common/include/kopano/ECGuid.h
	MUIDECSAB = DEFINE_GUID(0x50a921ac, 0xd340, 0x48ee, [8]byte{0xb3, 0x19, 0xfb, 0xa7, 0x53, 0x30, 0x44, 0x25})
)

type guidBytes struct {
	Data1 uint32
	Data2 uint16
	Data3 uint16
	Data4 [8]byte
}

// DEFINE_GUID is a helper to define byte representations of GUIDs.
func DEFINE_GUID(l uint32, w1, w2 uint16, b [8]byte) [16]byte {
	guid := guidBytes{
		l,
		w1,
		w2,
		b,
	}

	buf := bytes.NewBuffer(make([]byte, 0, 16))
	err := binary.Write(buf, binary.LittleEndian, guid)
	if err != nil {
		panic(err)
	}

	var res [16]byte
	copy(res[:], buf.Bytes())

	return res
}
