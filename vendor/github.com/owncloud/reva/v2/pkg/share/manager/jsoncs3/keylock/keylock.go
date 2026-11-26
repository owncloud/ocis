// Copyright 2026 ownCloud GmbH <mail@owncloud.com>
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

// Package keylock provides per-key mutual exclusion whose acquisition can be
// abandoned when the caller's context is done. The jsoncs3 share manager
// caches use it to guard their per-user/per-group/per-space shards, so that a
// request whose client already gave up does not stay parked (and keep its
// memory) behind a slow lock holder.
package keylock

import (
	"context"
	"runtime/debug"
	"sync"
	"sync/atomic"

	"github.com/owncloud/reva/v2/pkg/appctx"
)

// Map hands out one lock per key. The zero value is ready to use. Like the
// sync.Mutex-per-key map it replaces, entries are never removed, so the map
// grows with the number of distinct keys seen, not with the number of waiters.
type Map struct {
	locks sync.Map // key -> chan struct{} with capacity 1
}

// Lock acquires the lock for key, blocking until it is free or ctx is done.
// On success it returns the unlock function, which must be called exactly once;
// calling it again does nothing except log an error with a stack trace, using
// the logger from ctx.
// If ctx is done before or while waiting, it returns ctx.Err() and does not
// hold the lock.
func (m *Map) Lock(ctx context.Context, key string) (func(), error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	v, ok := m.locks.Load(key)
	if !ok {
		v, _ = m.locks.LoadOrStore(key, make(chan struct{}, 1))
	}
	sem := v.(chan struct{})

	select {
	case sem <- struct{}{}:
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	// select picks randomly when both cases are ready; never hand the lock to
	// a caller whose context is already done.
	if err := ctx.Err(); err != nil {
		<-sem
		return nil, err
	}
	// A second unlock is a bug. It is logged and otherwise ignored. The
	// per-holder flag also catches it while another goroutine already holds the
	// key, where draining the channel would silently release that goroutine's
	// lock.
	var released atomic.Bool
	return func() {
		if !released.CompareAndSwap(false, true) {
			logDoubleUnlock(ctx, key)
			return
		}
		select {
		case <-sem:
		default:
			// unreachable: only this holder can drain its own token
			logDoubleUnlock(ctx, key)
		}
	}, nil
}

func logDoubleUnlock(ctx context.Context, key string) {
	appctx.GetLogger(ctx).Error().
		Str("key", key).
		Str("stack", string(debug.Stack())).
		Msg("keylock: unlock of unlocked key")
}
