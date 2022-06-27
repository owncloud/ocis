/*
 * Copyright 2017-2019 Kopano
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

package oidc

import (
	"crypto"
	"encoding/base64"
	"fmt"
)

// LeftmostHashBytes defines []bytes with Base64URL encoder via String().
type LeftmostHashBytes []byte

// LeftmostHash hashes the provided data with the provided hash function and
// returns the left-most half the hashed bytes.
func LeftmostHash(data []byte, hash crypto.Hash) LeftmostHashBytes {
	hasher := hash.New()
	hasher.Write(data)
	result := hasher.Sum(nil)

	return LeftmostHashBytes(result[:len(result)/2])
}

// String returns the Base64URL encoded string of the accociated bytes.
func (lmhb LeftmostHashBytes) String() string {
	return base64.RawURLEncoding.EncodeToString(lmhb)
}

// HashFromSigningMethod returns the matching crypto.Hash for the provided
// signing alg.
func HashFromSigningMethod(alg string) (hash crypto.Hash, err error) {
	switch alg {
	case "HS256":
		fallthrough
	case "RS256":
		fallthrough
	case "PS256":
		fallthrough
	case "ES256":
		hash = crypto.SHA256

	case "HS386":
		fallthrough
	case "RS384":
		fallthrough
	case "PS384":
		fallthrough
	case "ES384":
		hash = crypto.SHA384

	case "HS512":
		fallthrough
	case "RS512":
		fallthrough
	case "PS512":
		fallthrough
	case "ES512":
		hash = crypto.SHA512

	case "EdDSA":
		hash = crypto.SHA512

	default:
		err = fmt.Errorf("Unkown alg %s", alg)
	}

	if !hash.Available() {
		err = fmt.Errorf("Hash for %s is not available", alg)
	}

	return
}
