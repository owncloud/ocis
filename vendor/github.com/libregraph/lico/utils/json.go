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
	"encoding/json"
	"net/http"
)

const (
	defaultJSONContentType = "application/json; encoding=utf-8"
)

// WriteJSON marshals the provided data as JSON and writes it to the provided
// http.ResponseWriter using the provided HTTP status code and content-type. the
// nature of this function is that it always writes a HTTP response header. Thus
// it makes no sense to write another header on error. Resulting errors should
// be logged and the connection should be closes as it is non-functional.
func WriteJSON(rw http.ResponseWriter, code int, data interface{}, contentType string) error {
	if contentType == "" {
		rw.Header().Set("Content-Type", defaultJSONContentType)
	} else {
		rw.Header().Set("content-Type", contentType)
	}

	rw.WriteHeader(code)

	enc := json.NewEncoder(rw)
	enc.SetIndent("", "  ")

	return enc.Encode(data)
}
