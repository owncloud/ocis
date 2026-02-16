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
	"fmt"
	"net/http"
	"net/url"
	"time"
)

var (
	farPastExpiryTime = time.Unix(0, 0)
)

func uniqueStrings(s []string) []string {
	var res []string
	m := make(map[string]bool)
	for _, s := range s {
		if _, ok := m[s]; ok {
			continue
		}
		m[s] = true
		res = append(res, s)
	}

	return res
}

func getRequestURL(req *http.Request, isTrustedSource bool) *url.URL {
	u, _ := url.Parse(req.URL.String())

	if isTrustedSource {
		// Look at proxy injected values to rewrite URLs if trusted.
		for {
			prefix := req.Header.Get("X-Forwarded-Prefix")
			if prefix != "" {
				u.Path = fmt.Sprintf("%s%s", prefix, u.Path)
				break
			}
			replaced := req.Header.Get("X-Replaced-Path")
			if replaced != "" {
				u.Path = replaced
				break
			}

			break
		}
	}

	return u
}

func addResponseHeaders(header http.Header) {
	header.Set("Cache-Control", "no-cache, no-store, must-revalidate")
	header.Set("Pragma", "no-cache")
	header.Set("X-Content-Type-Options", "nosniff")
	header.Set("Referrer-Policy", "origin")
}

func makeArrayFromBoolMap(m map[string]bool) []string {
	result := []string{}
	for k, v := range m {
		if v {
			result = append(result, k)
		}
	}

	return result
}
