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
	"fmt"
	"net/http"
)

// WriteErrorPage create a formatted error page response containing the provided
// information and writes it to the provided http.ResponseWriter.
func WriteErrorPage(rw http.ResponseWriter, code int, title string, message string) {
	if title == "" {
		title = http.StatusText(code)
	}

	text := fmt.Sprintf("%d %s", code, title)
	if message != "" {
		text = text + " - " + message
	}

	http.Error(rw, text, code)
}
