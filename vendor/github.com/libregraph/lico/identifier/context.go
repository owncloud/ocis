/*
 * Copyright 2021 Kopano and its licensors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package identifier

import (
	"context"

	"github.com/libregraph/lico/identifier/backends"
)

// Record is the struct which the identifier puts into the context.
type Record struct {
	HelloRequest    *HelloRequest
	UserFromBackend backends.UserFromBackend
}

// key is an unexported type for keys defined in this package.
// This prevents collisions with keys defined in other packages.
type key int

// recordKey is the key for identifier.Record in Contexts. It is
// unexported; clients use identifier.NewContext and identifier.FromContext
// instead of using this key directly.
var recordKey key

// NewRecordContext returns a new Context that carries value HelloRequest.
func NewRecordContext(ctx context.Context, record *Record) context.Context {
	return context.WithValue(ctx, recordKey, record)
}

// FromRecordContext returns the Record value stored in ctx, if any.
func FromRecordContext(ctx context.Context) (*Record, bool) {
	record, ok := ctx.Value(recordKey).(*Record)
	return record, ok
}
