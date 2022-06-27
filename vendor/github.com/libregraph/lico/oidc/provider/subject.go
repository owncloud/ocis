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
	"errors"

	konnect "github.com/libregraph/lico"
	"github.com/libregraph/lico/identity"
)

// PublicSubjectFromAuth creates the provideds auth Subject value with the
// accociated provider. This subject can be used as URL safe value to uniquely
// identify the provided auth user with remote systems.
func (p *Provider) PublicSubjectFromAuth(auth identity.AuthRecord) (string, error) {
	authorizedScopes := auth.AuthorizedScopes()
	if ok, _ := authorizedScopes[konnect.ScopeRawSubject]; ok {
		// Return raw subject as is when with ScopeRawSubject.
		user := auth.User()
		if user == nil {
			return "", errors.New("no user")
		}

		return user.Raw(), nil
	}

	return auth.Subject(), nil
}
