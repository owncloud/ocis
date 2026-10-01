// Package ratelimit implements a per-key sliding-window request limiter
// backed by a go-micro store, so the limit is correct even when the llm
// service runs multiple replicas. The algorithm mirrors reva's
// publicshares.BruteForceProtection (see
// vendor/github.com/owncloud/reva/v2/pkg/auth/manager/publicshares/bruteforceprotection.go),
// applied to generic per-user request counting instead of failed-auth
// attempts.
package ratelimit

import (
	"encoding/json"
	"errors"
	"sync"
	"time"

	microstore "go-micro.dev/v4/store"
)

// Limiter enforces a per-key sliding-window request limit.
type Limiter struct {
	store       microstore.Store
	window      time.Duration
	maxRequests int
	mutex       sync.Mutex
}

// New creates a Limiter. If window <= 0 or maxRequests <= 0, the limiter is
// disabled and Allow always returns true.
func New(store microstore.Store, window time.Duration, maxRequests int) *Limiter {
	return &Limiter{
		store:       store,
		window:      window,
		maxRequests: maxRequests,
	}
}

// Allow registers a request for key and reports whether it is within the
// configured limit. Expired timestamps are pruned on every call.
func (l *Limiter) Allow(key string) (bool, error) {
	if l == nil {
		return true, nil
	}
	if l.window <= 0 || l.maxRequests <= 0 {
		return true, nil
	}

	l.mutex.Lock()
	defer l.mutex.Unlock()

	timestamps, err := l.read(key)
	if err != nil {
		return false, err
	}

	cutoff := time.Now().Add(-l.window).UnixNano()
	fresh := make([]int64, 0, len(timestamps))
	for _, ts := range timestamps {
		if ts >= cutoff {
			fresh = append(fresh, ts)
		}
	}

	if len(fresh) >= l.maxRequests {
		return false, l.write(key, fresh)
	}

	fresh = append(fresh, time.Now().UnixNano())
	return true, l.write(key, fresh)
}

func (l *Limiter) read(key string) ([]int64, error) {
	records, err := l.store.Read(key)
	if errors.Is(err, microstore.ErrNotFound) {
		return []int64{}, nil
	}
	if err != nil {
		return nil, err
	}

	var timestamps []int64
	if err := json.Unmarshal(records[0].Value, &timestamps); err != nil {
		return nil, err
	}
	return timestamps, nil
}

func (l *Limiter) write(key string, timestamps []int64) error {
	if len(timestamps) == 0 {
		if err := l.store.Delete(key); err != nil && !errors.Is(err, microstore.ErrNotFound) {
			return err
		}
		return nil
	}

	value, err := json.Marshal(timestamps)
	if err != nil {
		return err
	}
	return l.store.Write(&microstore.Record{Key: key, Value: value})
}
