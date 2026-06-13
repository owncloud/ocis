package service

import (
	"encoding/json"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/owncloud/reva/v2/pkg/storagespace"
	"github.com/owncloud/reva/v2/pkg/store"
	"github.com/stretchr/testify/require"
	microstore "go-micro.dev/v4/store"

	"github.com/owncloud/ocis/v2/ocis-pkg/log"
)

// TestDebounceCoalescesActivities verifies that activities queued for the same
// resource within the window are flushed exactly once, as a single batch.
func TestDebounceCoalescesActivities(t *testing.T) {
	var (
		mu      sync.Mutex
		flushes int
		batched int
	)
	flush := func(_ string, acts []RawActivity) error {
		mu.Lock()
		defer mu.Unlock()
		flushes++
		batched += len(acts)
		return nil
	}

	d := NewDebouncer(100*time.Millisecond, log.NewLogger(), flush)
	for i := 0; i < 5; i++ {
		d.Debounce("resource-1", RawActivity{EventID: strconv.Itoa(i)})
	}

	// Nothing is flushed before the window elapses.
	mu.Lock()
	require.Equal(t, 0, flushes, "must not flush before the window elapses")
	mu.Unlock()

	require.Eventually(t, func() bool {
		mu.Lock()
		defer mu.Unlock()
		return flushes == 1
	}, time.Second, 5*time.Millisecond, "expected exactly one flush after the window")

	mu.Lock()
	require.Equal(t, 5, batched, "all queued activities must be flushed together")
	mu.Unlock()
}

// TestDebounceSynchronousWhenZero verifies that a zero window flushes immediately.
func TestDebounceSynchronousWhenZero(t *testing.T) {
	var got []RawActivity
	d := NewDebouncer(0, log.NewLogger(), func(_ string, acts []RawActivity) error {
		got = append(got, acts...)
		return nil
	})

	d.Debounce("resource-1", RawActivity{EventID: "a"})
	require.Len(t, got, 1, "zero window must flush synchronously")
	require.Equal(t, "a", got[0].EventID)
}

// TestActivitiesReadsLegacyJSON verifies the msgpack read path falls back to the
// previous json encoding so records written before the upgrade stay readable,
// and that appending re-encodes everything with msgpack.
func TestActivitiesReadsLegacyJSON(t *testing.T) {
	sto := store.Create()
	alog := &ActivitylogService{store: sto}
	alog.debouncer = NewDebouncer(0, log.NewLogger(), alog.storeActivity)

	rid := resourceID("legacy")
	key := storagespace.FormatResourceID(rid)

	// Seed a legacy json-encoded record directly.
	legacy, err := json.Marshal([]RawActivity{{EventID: "old", Depth: 1}})
	require.NoError(t, err)
	require.NoError(t, sto.Write(&microstore.Record{Key: key, Value: legacy}))

	got, err := alog.Activities(rid)
	require.NoError(t, err)
	require.Len(t, got, 1)
	require.Equal(t, "old", got[0].EventID)

	// Appending writes msgpack and must keep the legacy entry readable.
	require.NoError(t, alog.storeActivity(key, []RawActivity{{EventID: "new", Depth: 0}}))

	got, err = alog.Activities(rid)
	require.NoError(t, err)
	require.Len(t, got, 2)
	require.ElementsMatch(t, []string{"old", "new"}, []string{got[0].EventID, got[1].EventID})
}
