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

package payload

import (
	"encoding/json"
	"reflect"

	"github.com/gorilla/schema"
)

var decoder = schema.NewDecoder()
var encoder = schema.NewEncoder()

// DecodeSchema decodes request form data into the provided dst schema struct.
func DecodeSchema(dst interface{}, src map[string][]string) error {
	return decoder.Decode(dst, src)
}

// EncodeSchema encodes the provided src schema to the provided map.
func EncodeSchema(src interface{}, dst map[string][]string) error {
	return encoder.Encode(src, dst)
}

// ConvertOIDCClaimsRequest is a converter function for oidc.ClaimsRequest data
// provided in URL schema.
func ConvertOIDCClaimsRequest(value string) reflect.Value {
	v := ClaimsRequest{}

	if err := json.Unmarshal([]byte(value), &v); err != nil {
		return reflect.Value{}
	}

	return reflect.ValueOf(v)
}

func init() {
	decoder.IgnoreUnknownKeys(true)
	decoder.RegisterConverter(ClaimsRequest{}, ConvertOIDCClaimsRequest)
}
