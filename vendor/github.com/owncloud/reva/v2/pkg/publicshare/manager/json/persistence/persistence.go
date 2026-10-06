// Copyright 2018-2021 CERN
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// In applying this license, CERN does not waive the privileges and immunities
// granted to it by virtue of its status as an Intergovernmental Organization
// or submit itself to any jurisdiction.

package persistence

import "context"

// PublicShares is a map indexing publicshares by their ids
type PublicShares map[string]interface{}

// Persistence defines the interface for the json publicshare manager persistence layers
type Persistence interface {
	Init(context.Context) error
	// Read returns the share database for reading only. The result may alias
	// the implementation's own cache and be shared with concurrent readers, so
	// a caller that intends to mutate it - or to hand it to Write - must take
	// its own Copy first. Implementations must publish a new map rather than
	// mutate one they have already returned.
	Read(context.Context) (PublicShares, error)
	Write(context.Context, PublicShares) error
}

// Copy returns a copy of db that shares no mutable state with it: a fresh
// top-level map plus a fresh copy of each share's own map. Callers that mutate
// what Read returned need it, so that a reader still ranging over the same
// result after releasing its lock cannot race the mutation (see
// manager.UpdatePublicShare, which rewrites an existing share's fields).
func Copy(db PublicShares) PublicShares {
	out := make(PublicShares, len(db))
	for k, v := range db {
		if entry, ok := v.(map[string]interface{}); ok {
			entryCopy := make(map[string]interface{}, len(entry))
			for ek, ev := range entry {
				entryCopy[ek] = ev
			}
			out[k] = entryCopy
			continue
		}
		out[k] = v
	}
	return out
}
