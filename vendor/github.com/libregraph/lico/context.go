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

package lico

import (
	"context"

	"github.com/golang-jwt/jwt/v4"
)

// key is an unexported type for keys defined in this package.
// This prevents collisions with keys defined in other packages.
type key int

// claimsKey is the key for claims in contexts. It is
// unexported; clients use konnect.NewClaimsContext and
// connect.FromClaimsContext instead of using this key directly.
var claimsKey key

// NewClaimsContext returns a new Context that carries value auth.
func NewClaimsContext(ctx context.Context, claims jwt.Claims) context.Context {
	return context.WithValue(ctx, claimsKey, claims)
}

// FromClaimsContext returns the AuthRecord value stored in ctx, if any.
func FromClaimsContext(ctx context.Context) (jwt.Claims, bool) {
	claims, ok := ctx.Value(claimsKey).(jwt.Claims)
	return claims, ok
}
