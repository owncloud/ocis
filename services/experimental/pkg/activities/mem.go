package activities

import (
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

// MemStorage holds activity data
type MemStorage struct {
	capacity, length uint64
	data             sync.Map
}

// NewMemStore bootstraps memory store.
func NewMemStore(capacity uint64) *MemStorage {
	return &MemStorage{
		capacity: capacity,
	}
}

// List loads all entries
func (s *MemStorage) List(uID string) []Activity {
	o := make([]Activity, 0)
	for _, k := range s.order() {
		if v, ok := s.data.Load(k); ok {
			a := v.(Activity)

			if uID != a.UserID {
				continue
			}

			o = append(o, v.(Activity))
		}
	}

	return o
}

// Add adds an entry for given key and value
func (s *MemStorage) Add(a Activity) {
	s.data.Store(time.Now(), a)
	atomic.AddUint64(&s.length, 1)

	if s.length > s.capacity {
		for _, k := range s.order()[s.capacity:] {
			s.data.Delete(k)
			atomic.AddUint64(&s.length, ^uint64(0))
		}
	}
}

func (s *MemStorage) order() (o []time.Time) {
	var keys []time.Time
	s.data.Range(func(k, _ any) bool {
		keys = append(keys, k.(time.Time))
		return true
	})

	sort.Slice(keys, func(i, j int) bool {
		return keys[i].After(keys[j])
	})

	for _, k := range keys {
		o = append(o, k)
	}

	return
}
