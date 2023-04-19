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

package code

import (
	"github.com/libregraph/lico/identity"
	"github.com/libregraph/lico/oidc/payload"
)

// Record bundles the data storedi in a code manager.
type Record struct {
	AuthenticationRequest *payload.AuthenticationRequest
	Auth                  identity.AuthRecord
	Session               *payload.Session
}

// Manager is a interface defining a code manager.
type Manager interface {
	Create(record *Record) (string, error)
	Pop(code string) (*Record, bool)
}
