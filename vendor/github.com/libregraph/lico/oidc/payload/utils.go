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
)

// ToMap is a helper function to convert the provided payload struct to
// a map type which can be used to extend the payload data with additional fields.
func ToMap(payload interface{}) (map[string]interface{}, error) {
	// NOTE(longsleep): This implementation sucks, marshal to JSON and unmarshal
	// again - rly?
	intermediate, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	claims := make(map[string]interface{})
	err = json.Unmarshal(intermediate, &claims)
	if err != nil {
		return nil, err
	}

	return claims, nil
}
