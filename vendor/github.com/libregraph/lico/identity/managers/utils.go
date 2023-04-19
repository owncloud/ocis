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

package managers

import (
	"encoding/base64"

	"golang.org/x/crypto/blake2b"

	konnectoidc "github.com/libregraph/lico/oidc"
)

func setupSupportedScopes(scopes []string, extra []string, override []string) []string {
	if len(override) > 0 {
		return override
	}

	return append(scopes, extra...)
}

func getPublicSubject(sub []byte, extra []byte) (string, error) {
	// Hash the raw subject with our specific salt.
	hasher, err := blake2b.New512([]byte(konnectoidc.LibreGraphIDTokenSubjectSaltV1))
	if err != nil {
		return "", err
	}

	hasher.Write(sub)
	hasher.Write([]byte(" "))
	hasher.Write(extra)

	// NOTE(longsleep): URL safe encoding for subject is important since many
	// third party applications validate this with rather strict patterns. We
	// also inject an @ to ensure its compatible to some apps which require one.
	s := base64.RawURLEncoding.EncodeToString(hasher.Sum(nil))
	return s[:16] + "@" + s[16:], nil
}
