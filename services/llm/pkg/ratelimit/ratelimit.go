// Package ratelimit implements a per-key sliding-window request limiter
// backed by a go-micro store.
//
// The go-micro store interface has no compare-and-swap, so a read-modify-write
// of a single shared counter would let concurrent replicas overwrite each
// other. Instead every request writes its own record (key "<key>/<unixnano>-<nonce>")
// and then counts the key's records within the window. A request that finds
// itself over the limit removes its own record and is denied. Because records
// are never overwritten, the limit holds across replicas: concurrent requests
// at the boundary may all be denied, but more than the maximum are never
// admitted.
//
// Records are expired by the store (a per-write TTL where supported, the
// store-wide TTL otherwise, which must therefore be at least the window).
package ratelimit

import (
	"crypto/rand"
	"encoding/hex"
	"strconv"
	"strings"
	"time"

	microstore "go-micro.dev/v4/store"
)

const separator = "/"

// Limiter enforces a per-key sliding-window request limit.
type Limiter struct {
	store       microstore.Store
	window      time.Duration
	maxRequests int
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
// configured limit. A denied request does not count against the limit.
func (l *Limiter) Allow(key string) (bool, error) {
	if l == nil {
		return true, nil
	}
	if l.window <= 0 || l.maxRequests <= 0 {
		return true, nil
	}

	nonce := make([]byte, 4)
	if _, err := rand.Read(nonce); err != nil {
		return false, err
	}
	now := time.Now()
	recordKey := key + separator + strconv.FormatInt(now.UnixNano(), 10) + "-" + hex.EncodeToString(nonce)

	if err := l.store.Write(&microstore.Record{Key: recordKey}, microstore.WriteTTL(l.window)); err != nil {
		return false, err
	}

	prefix := key + separator
	keys, err := l.store.List(microstore.ListPrefix(prefix))
	if err != nil {
		// do not leave an uncounted-for record behind when we cannot decide
		_ = l.store.Delete(recordKey)
		return false, err
	}

	cutoff := now.Add(-l.window).UnixNano()
	count := 0
	for _, k := range keys {
		if recordTime(k, prefix) >= cutoff {
			count++
		}
	}

	if count > l.maxRequests {
		if err := l.store.Delete(recordKey); err != nil && err != microstore.ErrNotFound {
			return false, err
		}
		return false, nil
	}
	return true, nil
}

// recordTime extracts the request timestamp from a record key. Keys that do
// not belong to prefix or are malformed yield 0 (i.e. "expired").
func recordTime(recordKey, prefix string) int64 {
	if !strings.HasPrefix(recordKey, prefix) {
		return 0
	}
	suffix := strings.TrimPrefix(recordKey, prefix)
	ts, _, _ := strings.Cut(suffix, "-")
	n, err := strconv.ParseInt(ts, 10, 64)
	if err != nil {
		return 0
	}
	return n
}
