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
	"github.com/gorilla/schema"
)

// Create a Decoder instance as a package global, because it caches
// meta-data about structs, and an instance can be shared safely.
var urlSchemaDecoder = schema.NewDecoder()

// DecodeURLSchema decodes request for mdata in to the provided dst url struct.
func DecodeURLSchema(dst interface{}, src map[string][]string) error {
	return urlSchemaDecoder.Decode(dst, src)
}

func init() {
	urlSchemaDecoder.SetAliasTag("url")
	urlSchemaDecoder.IgnoreUnknownKeys(true)
}
