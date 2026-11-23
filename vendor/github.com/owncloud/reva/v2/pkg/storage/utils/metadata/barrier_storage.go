package metadata

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// barrierTimeout bounds how long an arrived Upload call waits for the
// remaining callers before giving up, so a caller that never arrives (e.g.
// because its own cold-start sync errored before reaching Upload) can't
// hang the others for the full test-binary timeout.
const barrierTimeout = 2 * time.Second

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
	case <-time.After(barrierTimeout):
		return nil, fmt.Errorf("barrier: timed out after %s waiting for %d/%d callers to arrive", barrierTimeout, atomic.LoadInt32(&b.arrived), b.n)
	}
	return b.Storage.Upload(ctx, req)
}

// SimpleUpload routes through Upload so it also counts toward and waits on
// the barrier, instead of promoting straight to the embedded Storage.
func (b *BarrierStorage) SimpleUpload(ctx context.Context, uploadpath string, content []byte) error {
	_, err := b.Upload(ctx, UploadRequest{
		Path:    uploadpath,
		Content: content,
	})
	return err
}
