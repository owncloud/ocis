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

package cs3

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/owncloud/reva/v2/pkg/errtypes"
	"github.com/owncloud/reva/v2/pkg/publicshare/manager/json/persistence"
	"github.com/owncloud/reva/v2/pkg/storage/utils/metadata"
	"github.com/owncloud/reva/v2/pkg/utils"
)

type db struct {
	mtime        time.Time
	publicShares persistence.PublicShares
}

type cs3 struct {
	// initialized is read on every call (Init is called unconditionally by
	// the manager ahead of every Read/Write) and must stay lock-free once
	// true, or every already-initialized caller would queue behind mu for a
	// single bool check - and worse, behind whatever Read is doing with mu
	// at the time, since they'd be contending for the same lock.
	initialized atomic.Bool
	s           metadata.Storage
	// initMu serializes the one-time (and retried-on-failure) call into
	// s.Init. It is deliberately not mu: that state is independent of the
	// mtime cache below, and sharing a lock would reintroduce the same
	// contention initialized/atomic.Bool exists to avoid.
	initMu sync.Mutex
	// cs3 caches the remote publicshares.json in p.db and only refetches it when
	// its mtime advances. That cache refill mutates p.db as a side effect of
	// Read, so Read cannot be treated as a pure/reentrant read by callers - mu
	// serializes access to p.db regardless of how the caller itself locks.
	mu sync.Mutex
	db db
}

// New returns a new Cache instance
func New(s metadata.Storage) persistence.Persistence {
	return &cs3{
		s: s,
		db: db{
			publicShares: persistence.PublicShares{},
		},
	}
}

func (p *cs3) Init(ctx context.Context) error {
	if p.initialized.Load() {
		return nil
	}

	p.initMu.Lock()
	defer p.initMu.Unlock()

	if p.initialized.Load() {
		return nil
	}

	err := p.s.Init(ctx, "jsoncs3-public-share-manager-metadata")
	if err != nil {
		return err
	}
	p.initialized.Store(true)

	return nil
}

func (p *cs3) Read(ctx context.Context) (persistence.PublicShares, error) {
	if !p.initialized.Load() {
		return nil, fmt.Errorf("not initialized")
	}

	// We use the Lock because the read function updates the cache. So most time operations should be fast.
	p.mu.Lock()
	defer p.mu.Unlock()

	info, err := p.s.Stat(ctx, "publicshares.json")
	if err != nil {
		if _, ok := err.(errtypes.NotFound); ok {
			return persistence.Copy(p.db.publicShares), nil // Nothing to sync against
		}
		return nil, err
	}

	if utils.TSToTime(info.Mtime).After(p.db.mtime) {
		readBytes, err := p.s.SimpleDownload(ctx, "publicshares.json")
		if err != nil {
			return nil, err
		}
		p.db.publicShares = persistence.PublicShares{}
		if err := json.Unmarshal(readBytes, &p.db.publicShares); err != nil {
			return nil, err
		}
		p.db.mtime = utils.TSToTime(info.Mtime)
	}
	return persistence.Copy(p.db.publicShares), nil
}

func (p *cs3) Write(ctx context.Context, db persistence.PublicShares) error {
	if !p.initialized.Load() {
		return fmt.Errorf("not initialized")
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	dbAsJSON, err := json.Marshal(db)
	if err != nil {
		return err
	}

	_, err = p.s.Upload(ctx, metadata.UploadRequest{
		Content:           dbAsJSON,
		Path:              "publicshares.json",
		IfUnmodifiedSince: p.db.mtime,
	})
	if err != nil {
		return err
	}

	// Keep the cache in sync with what was just persisted. This used to
	// happen implicitly, because Read() handed out a reference to
	// p.db.publicShares itself and callers mutated it in place before
	// calling Write() with that same map. Now that Read() returns an
	// independent copy (see persistence.Copy), it has to be done explicitly
	// here, or the cache would only pick up our own write once some later
	// external write advances the remote mtime past our stale one.
	if info, statErr := p.s.Stat(ctx, "publicshares.json"); statErr == nil {
		p.db.mtime = utils.TSToTime(info.Mtime)
	}
	p.db.publicShares = persistence.Copy(db)

	return nil
}
