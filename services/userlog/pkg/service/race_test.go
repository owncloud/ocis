package service

import (
	"context"
	"encoding/json"
	"strconv"
	"sync"
	"testing"
	"time"

	microstore "go-micro.dev/v4/store"
	"go.opentelemetry.io/otel/trace"
)

// TestAlterUserEventListConcurrentDoesNotLoseEvents asserts that concurrent
// modifications of a single user's event list do not drop updates. Events are
// processed by MaxConcurrency workers, so two of them appending to the same
// user list must not lose each other's writes. The small sleep widens the
// read-modify-write window so the missing-lock case fails deterministically.
func TestAlterUserEventListConcurrentDoesNotLoseEvents(t *testing.T) {
	ul := &UserlogService{store: microstore.NewMemoryStore()}

	const userID = "user-1"
	const n = 50

	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			id := strconv.Itoa(i)
			if err := ul.alterUserEventList(userID, func(ids []string) []string {
				time.Sleep(time.Millisecond)
				return append(ids, id)
			}); err != nil {
				t.Errorf("alterUserEventList: %v", err)
			}
		}(i)
	}
	wg.Wait()

	recs, err := ul.store.Read(userID)
	if err != nil {
		t.Fatalf("read back user event list: %v", err)
	}
	var ids []string
	if len(recs) > 0 {
		if err := json.Unmarshal(recs[0].Value, &ids); err != nil {
			t.Fatalf("unmarshal user event list: %v", err)
		}
	}
	if len(ids) != n {
		t.Fatalf("lost user events under concurrency: got %d, want %d", len(ids), n)
	}
}

// TestAlterGlobalEventsConcurrentDoesNotLoseEvents asserts the same atomicity
// for the single shared global-events key, which every global event mutates.
func TestAlterGlobalEventsConcurrentDoesNotLoseEvents(t *testing.T) {
	ul := &UserlogService{
		store:  microstore.NewMemoryStore(),
		tracer: trace.NewNoopTracerProvider().Tracer(""),
	}

	const n = 50

	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := strconv.Itoa(i)
			if err := ul.alterGlobalEvents(context.Background(), func(evs map[string]json.RawMessage) error {
				time.Sleep(time.Millisecond)
				evs[key] = json.RawMessage(`"x"`)
				return nil
			}); err != nil {
				t.Errorf("alterGlobalEvents: %v", err)
			}
		}(i)
	}
	wg.Wait()

	evs, err := ul.GetGlobalEvents(context.Background())
	if err != nil {
		t.Fatalf("GetGlobalEvents: %v", err)
	}
	if len(evs) != n {
		t.Fatalf("lost global events under concurrency: got %d, want %d", len(evs), n)
	}
}
