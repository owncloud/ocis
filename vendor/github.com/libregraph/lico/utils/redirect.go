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

package utils

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/google/go-querystring/query"
)

// WriteRedirect crates a URL out of the provided uri and params and writes a
// a HTTP response with the provided HTTP status code to the provided
// http.ResponseWriter incliding HTTP caching headers to prevent caching. If
// asFragment is true, the provided params are added as URL fragment, otherwise
// they replace the query. If params is nil, the provided uri is taken as is.
func WriteRedirect(rw http.ResponseWriter, code int, uri *url.URL, params interface{}, asFragment bool) error {
	if params != nil {
		paramValues, err := query.Values(params)
		if err != nil {
			return err
		}

		u, _ := url.Parse(uri.String())
		if asFragment {
			if u.Fragment != "" {
				u.Fragment += "&"
			}
			f := paramValues.Encode()   // This encods into URL encoded form with QueryEscape.
			f, _ = url.QueryUnescape(f) // But we need it unencoded since its the fragment, it is encoded later (when serializing the URL).
			u.Fragment += f             // Append fragment extension.
		} else {
			queryValues := u.Query()
			for k, vs := range paramValues {
				for _, v := range vs {
					queryValues.Add(k, v)
				}
			}
			u.RawQuery = strings.ReplaceAll(queryValues.Encode(), "+", "%20") // NOTE(longsleep): Ensure we use %20 instead of +.
		}
		uri = u
	}

	rw.Header().Set("Location", uri.String())
	rw.Header().Set("Cache-Control", "no-store")
	rw.Header().Set("Pragma", "no-cache")

	rw.WriteHeader(code)

	return nil
}
