// Copyright 2022 CERN
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

package ldap

import (
	"crypto/tls"
	"time"
)

// Config holds the basic configuration of the LDAP Connection
type Config struct {
	URI          string
	BindDN       string
	BindPassword string
	TLSConfig    *tls.Config

	RetryMaxCount  int
	RetryBaseDelay time.Duration
	RetryMaxDelay  time.Duration

	// PoolSize caps the number of concurrently open connections in the pool. Only used by
	// NewLDAPPool; NewLDAPWithReconnect ignores it. Defaults to defaultPoolSize (5) when <= 0.
	PoolSize int
	// PoolCheckoutTimeout bounds how long a checkout blocks once the pool is at PoolSize. Only used
	// by NewLDAPPool. Defaults to defaultPoolCheckoutTimeout (30s) when <= 0.
	PoolCheckoutTimeout time.Duration
}
