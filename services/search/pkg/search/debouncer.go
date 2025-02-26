package search

import (
	"sync"
	"time"

	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
)

// SpaceDebouncer debounces operations on spaces for a configurable amount of time
type SpaceDebouncer struct {
	after      time.Duration
	f          func(id *provider.StorageSpaceId)
	pending    map[string]*time.Timer
	inProgress sync.Map

	mutex sync.Mutex
}

// NewSpaceDebouncer returns a new SpaceDebouncer instance
func NewSpaceDebouncer(d time.Duration, f func(id *provider.StorageSpaceId)) *SpaceDebouncer {
	return &SpaceDebouncer{
		after:      d,
		f:          f,
		pending:    map[string]*time.Timer{},
		inProgress: sync.Map{},
	}
}

// Debounce restars the debounce timer for the given space
func (d *SpaceDebouncer) Debounce(id *provider.StorageSpaceId) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if t := d.pending[id.OpaqueId]; t != nil {
		t.Stop()
	}

	d.pending[id.OpaqueId] = time.AfterFunc(d.after, func() {
		if _, ok := d.inProgress.Load(id.OpaqueId); ok {
			// Reschedule this run for when the previous run has finished
			d.mutex.Lock()
			d.pending[id.OpaqueId].Reset(d.after)
			d.mutex.Unlock()
			return
		}

		d.inProgress.Store(id.OpaqueId, true)
		defer d.inProgress.Delete(id.OpaqueId)
		d.f(id)
	})
}
