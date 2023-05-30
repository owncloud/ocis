// Copyright 2018-2021 CERN
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// In applying this license, CERN does not waive the privileges and immunities
// granted to it by virtue of its status as an Intergovernmental Organization
// or submit itself to any jurisdiction.

package pool

import (
	"sync"
)

// TODO(labkode): is concurrent access to the maps safe?
// var storageProviders = map[string]storageprovider.ProviderAPIClient{}
var (
	storageProviders       = newProvider()
	authProviders          = newProvider()
	appAuthProviders       = newProvider()
	authRegistries         = newProvider()
	userShareProviders     = newProvider()
	ocmShareProviders      = newProvider()
	ocmInviteManagers      = newProvider()
	ocmProviderAuthorizers = newProvider()
	ocmCores               = newProvider()
	publicShareProviders   = newProvider()
	preferencesProviders   = newProvider()
	permissionsProviders   = newProvider()
	appRegistries          = newProvider()
	appProviders           = newProvider()
	storageRegistries      = newProvider()
	gatewayProviders       = newProvider()
	userProviders          = newProvider()
	groupProviders         = newProvider()
	dataTxs                = newProvider()
	maxCallRecvMsgSize     = 10240000
)

type provider struct {
	m    sync.Mutex
	conn map[string]interface{}
}

func newProvider() provider {
	return provider{
		sync.Mutex{},
		make(map[string]interface{}),
	}
}
