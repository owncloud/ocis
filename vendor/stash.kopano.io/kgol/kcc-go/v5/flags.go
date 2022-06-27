/*
 * Copyright 2017-2019 Kopano and its licensors
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
	"strconv"
)

// KCFlag is the type representing flags known to Kopano Core.
type KCFlag uint64

func (kcf KCFlag) String() string {
	return strconv.FormatUint(uint64(kcf), 10)
}

// ABFlag is the type representing flags known to Kopano Core for AB.
type ABFlag uint64

func (abf ABFlag) String() string {
	return strconv.FormatUint(uint64(abf), 10)
}

// Kopano capability flags as defined in provider/include/kcore.hpp. This only
// defines the flags actually used or understood by kcc-go.
const (
	KOPANO_CAP_LARGE_SESSIONID KCFlag = 0x0010
	KOPANO_CAP_MULTI_SERVER    KCFlag = 0x0040
	KOPANO_CAP_ENHANCED_ICS    KCFlag = 0x0100
	KOPANO_CAP_UNICODE         KCFlag = 0x0200
)

// DefaultClientCapabilities groups the default client caps sent by kcc.
var DefaultClientCapabilities = KOPANO_CAP_UNICODE |
	KOPANO_CAP_LARGE_SESSIONID |
	KOPANO_CAP_MULTI_SERVER |
	KOPANO_CAP_ENHANCED_ICS

// Kopano logon flags as defined in provider/include/kcore.hpp. This only
// defines the flags actually used or understood by kcc-go.
const (
	KOPANO_LOGON_NO_UID_AUTH         KCFlag = 0x0001
	KOPANO_LOGON_NO_REGISTER_SESSION KCFlag = 0x0002
)

// Kopano AB flags as defined in mapi4linux/include/mapidefs.h. This only
// defines the flags actually used or understood by kcc-go.
const (
	MAPI_UNRESOLVED ABFlag = 0x00000000
	MAPI_AMBIGUOUS  ABFlag = 0x00000001
	MAPI_RESOLVED   ABFlag = 0x00000002
)
