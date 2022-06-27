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

package bsguest

import (
	"github.com/libregraph/lico/bootstrap"
	"github.com/libregraph/lico/identity"
	"github.com/libregraph/lico/identity/managers"
)

// Identity managers.
const (
	identityManagerName = "guest"
)

func Register() error {
	return bootstrap.RegisterIdentityManager(identityManagerName, NewIdentityManager)
}

func MustRegister() {
	if err := Register(); err != nil {
		panic(err)
	}
}

func NewIdentityManager(bs bootstrap.Bootstrap) (identity.Manager, error) {
	config := bs.Config()

	logger := config.Config.Logger

	identityManagerConfig := &identity.Config{
		Logger: logger,
	}

	guestIdentityManager := managers.NewGuestIdentityManager(identityManagerConfig)

	return guestIdentityManager, nil
}
