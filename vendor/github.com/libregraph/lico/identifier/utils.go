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

package identifier

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gabriel-vasile/mimetype"
)

var (
	farPastExpiryTime                 = time.Unix(0, 0)
	farPastExpiryTimeHTTPHeaderString = farPastExpiryTime.UTC().Format(http.TimeFormat)
)

func addCommonResponseHeaders(header http.Header) {
	header.Set("X-Frame-Options", "DENY")
	header.Set("X-XSS-Protection", "1; mode=block")
	header.Set("X-Content-Type-Options", "nosniff")
	header.Set("Referrer-Policy", "origin")
}

func addNoCacheResponseHeaders(header http.Header) {
	header.Set("Cache-Control", "no-cache, no-store, must-revalidate")
	header.Set("Pragma", "no-cache")
	header.Set("Expires", farPastExpiryTimeHTTPHeaderString)
}

func encodeImageAsDataURL(b []byte) (string, error) {
	mt := mimetype.Detect(b)
	if !strings.HasPrefix(mt.String(), "image/") {
		return "", fmt.Errorf("not an image: %s", mt)
	}

	return "data:" + mt.String() + ";base64," + base64.StdEncoding.EncodeToString(b), nil
}
