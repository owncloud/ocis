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
	"strconv"
)

// MAPIType is the type representing MAPI types as used by Kopano Core.
type MAPIType uint32

func (mt MAPIType) String() string {
	return strconv.FormatUint(uint64(mt), 10)
}

// Possible type values as defined in mapi4linux/include/mapidefs.h. We
// only define the ones know and understood by kcc-go.
const (
	MAPI_MAILUSER MAPIType = 0x00000006
)
