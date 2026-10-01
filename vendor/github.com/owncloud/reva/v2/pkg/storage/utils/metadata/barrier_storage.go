package metadata

import (
	"context"
	"sync"
	"sync/atomic"
)

// BarrierStorage wraps a Storage and holds Upload calls until n goroutines have
// arrived, then releases them all at once, to deterministically reproduce
// concurrent-write races in tests.
type BarrierStorage struct {
	Storage
	arrived   int32
	n         int32
	ready     chan struct{}
	closeOnce sync.Once
}

// NewBarrierStorage returns a BarrierStorage wrapping s that releases once n Upload calls have arrived.
func NewBarrierStorage(s Storage, n int) *BarrierStorage {
	return &BarrierStorage{Storage: s, n: int32(n), ready: make(chan struct{})}
}

// Upload blocks until n calls have arrived, then delegates to the wrapped Storage.
func (b *BarrierStorage) Upload(ctx context.Context, req UploadRequest) (*UploadResponse, error) {
	if atomic.AddInt32(&b.arrived, 1) >= b.n {
		b.closeOnce.Do(func() { close(b.ready) })
	}
	select {
	case <-b.ready:
	case <-ctx.Done():
		return nil, ctx.Err()
	}
	return b.Storage.Upload(ctx, req)
}
