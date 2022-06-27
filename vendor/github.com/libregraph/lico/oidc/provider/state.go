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

package provider

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"

	"golang.org/x/crypto/blake2b"
	"stash.kopano.io/kgol/rndm"

	"github.com/libregraph/lico/identity"
	"github.com/libregraph/lico/oidc/payload"
)

var browserStateMarker = []byte("kopano-konnect-1")

func (p *Provider) makeBrowserState(ar *payload.AuthenticationRequest, auth identity.AuthRecord, err error) (string, error) {
	hasher, hasherErr := blake2b.New256(nil)
	if hasherErr != nil {
		return "", hasherErr
	}
	if auth != nil && err == nil {
		hasher.Write([]byte(auth.Subject()))
	} else {
		// Use empty string value when not signed in or with error. This means
		// that a browser state is always created.
		hasher.Write([]byte(" "))
	}
	hasher.Write([]byte(" "))
	hasher.Write([]byte(p.issuerIdentifier))
	hasher.Write([]byte(" "))
	hasher.Write(browserStateMarker)

	// Encode to string.
	browserState := base64.RawURLEncoding.EncodeToString(hasher.Sum(nil))

	return browserState, nil
}

func (p *Provider) makeSessionState(req *http.Request, ar *payload.AuthenticationRequest, browserState string) (string, error) {
	var origin string

	for {
		redirectURL := ar.RedirectURI
		if redirectURL != nil {
			origin = fmt.Sprintf("%s://%s", redirectURL.Scheme, redirectURL.Host)
			break
		}

		originHeader := req.Header.Get("Origin")
		if originHeader != "" {
			origin = originHeader
			break
		}

		refererHeader := req.Header.Get("Referer")
		if refererHeader != "" {
			// Rescure from referer.
			refererURL, err := url.Parse(refererHeader)
			if err != nil {
				return "", fmt.Errorf("invalid referer value: %v", err)
			}

			origin = fmt.Sprintf("%s://%s", refererURL.Scheme, refererURL.Host)
			break
		}

		return "", fmt.Errorf("missing origin")
	}

	salt := rndm.GenerateRandomString(32)

	hasher := sha256.New()
	hasher.Write([]byte(ar.ClientID))
	hasher.Write([]byte(" "))
	hasher.Write([]byte(origin))
	hasher.Write([]byte(" "))
	hasher.Write([]byte(browserState))
	hasher.Write([]byte(" "))
	hasher.Write([]byte(salt))

	sessionState := fmt.Sprintf("%s.%s", hex.EncodeToString(hasher.Sum(nil)), salt)

	return sessionState, nil
}
