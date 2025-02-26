/*
 * Copyright 2017-2019 Kopano and its licensors
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

package identity

import (
	"context"
)

// key is an unexported type for keys defined in this package.
// This prevents collisions with keys defined in other packages.
type key int

// authRecordKey is the key for identity.AuthRecord in Contexts. It is
// unexported; clients use identity.NewContext and identity.FromContext
// instead of using this key directly.
const (
	authRecordKey key = iota
)

// NewContext returns a new Context that carries value auth.
func NewContext(ctx context.Context, auth AuthRecord) context.Context {
	return context.WithValue(ctx, authRecordKey, auth)
}

// FromContext returns the AuthRecord value stored in ctx, if any.
func FromContext(ctx context.Context) (AuthRecord, bool) {
	auth, ok := ctx.Value(authRecordKey).(AuthRecord)
	return auth, ok
}
