/*
 * Copyright 2021 Kopano and its licensors
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
	"context"
	"net"
	"net/http"

	konnect "github.com/libregraph/lico"
	"github.com/libregraph/lico/config"
	"github.com/libregraph/lico/identifier/backends"
	"github.com/libregraph/lico/utils"
)

// Record is the struct which the identifier puts into the context.
type Record struct {
	HelloRequest *HelloRequest
	RealIP       string
	UserAgent    string

	BackendUser    backends.UserFromBackend
	IdentifiedUser *IdentifiedUser
}

func NewRecord(req *http.Request, c *config.Config) *Record {
	record := &Record{
		UserAgent: req.UserAgent(),
	}

	trusted, _ := utils.IsRequestFromTrustedSource(req, c.TrustedProxyIPs, c.TrustedProxyNets)

	if trusted {
		record.RealIP = req.Header.Get("X-Real-Ip")
	}
	if record.RealIP == "" {
		if ip, _, err := net.SplitHostPort(req.RemoteAddr); err == nil {
			record.RealIP = ip
		}
	}

	return record
}

// key is an unexported type for keys defined in this package.
// This prevents collisions with keys defined in other packages.
type key int

// Keys for context data.
// Unexported; Clients use identifier.New{?}Context and
// identifier.From{?}Context functions instead of using these keys directly.
const (
	recordKey key = iota
)

// NewRecordContext returns a new Context that carries the Record.
func NewRecordContext(ctx context.Context, record *Record) context.Context {
	return context.WithValue(ctx, recordKey, record)
}

// FromRecordContext returns the Record value stored in ctx, if any.
func FromRecordContext(ctx context.Context) (*Record, bool) {
	record, ok := ctx.Value(recordKey).(*Record)
	return record, ok
}

// RecordFromRequestContext returns a new Record value based on the request
// stored in ctx, if any.
func RecordFromRequestContext(ctx context.Context, c *config.Config) (*Record, bool) {
	if req, ok := konnect.FromRequestContext(ctx); ok {
		return NewRecord(req, c), true
	}
	return nil, false
}
