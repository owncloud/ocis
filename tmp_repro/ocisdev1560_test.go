package repro

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"testing"
	"time"

	collaboration "github.com/cs3org/go-cs3apis/cs3/sharing/collaboration/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/owncloud/reva/v2/pkg/errtypes"
	"github.com/owncloud/reva/v2/pkg/share/manager/jsoncs3/providercache"
	"github.com/owncloud/reva/v2/pkg/share/manager/jsoncs3/receivedsharecache"
	"github.com/owncloud/reva/v2/pkg/storage/utils/metadata"
)

// fakeStorage simulates storage-system: every Upload/Download is a network
// round trip (sleeps uploadLatency/downloadLatency), and Upload keeps
// returning errtypes.Aborted (etag mismatch) for the first failUploads
// attempts, modeling the etag-retry storm from OCISDEV-1560.
type fakeStorage struct {
	mu sync.Mutex

	uploadLatency   time.Duration
	downloadLatency time.Duration
	failUploads     int

	uploadCalls int
	etag        string
	content     []byte

	onFirstDownload func()
	firstDownload   sync.Once
}

func (f *fakeStorage) Backend() string                               { return "fake" }
func (f *fakeStorage) Init(ctx context.Context, name string) error   { return nil }
func (f *fakeStorage) SimpleUpload(ctx context.Context, p string, c []byte) error {
	return nil
}
func (f *fakeStorage) SimpleDownload(ctx context.Context, p string) ([]byte, error) {
	return nil, nil
}
func (f *fakeStorage) Delete(ctx context.Context, p string) error { return nil }
func (f *fakeStorage) Stat(ctx context.Context, p string) (*provider.ResourceInfo, error) {
	return nil, nil
}
func (f *fakeStorage) ReadDir(ctx context.Context, p string) ([]string, error) { return nil, nil }
func (f *fakeStorage) ListDir(ctx context.Context, p string) ([]*provider.ResourceInfo, error) {
	return nil, nil
}
func (f *fakeStorage) CreateSymlink(ctx context.Context, oldname, newname string) error {
	return nil
}
func (f *fakeStorage) ResolveSymlink(ctx context.Context, name string) (string, error) {
	return "", nil
}
func (f *fakeStorage) MakeDirIfNotExist(ctx context.Context, name string) error { return nil }
func (f *fakeStorage) WasRecentlyDeleted(ctx context.Context, p string) (bool, error) {
	return false, nil
}

func (f *fakeStorage) Download(ctx context.Context, req metadata.DownloadRequest) (*metadata.DownloadResponse, error) {
	f.firstDownload.Do(func() {
		if f.onFirstDownload != nil {
			f.onFirstDownload() // fires the instant the writer has grabbed LockSpace
		}
	})
	time.Sleep(f.downloadLatency)

	f.mu.Lock()
	defer f.mu.Unlock()
	if f.content == nil {
		return nil, errtypes.NotFound("space.json")
	}
	return &metadata.DownloadResponse{Content: f.content, Etag: f.etag}, nil
}

func (f *fakeStorage) Upload(ctx context.Context, req metadata.UploadRequest) (*metadata.UploadResponse, error) {
	time.Sleep(f.uploadLatency)

	f.mu.Lock()
	defer f.mu.Unlock()
	f.uploadCalls++
	if f.uploadCalls <= f.failUploads {
		return nil, errtypes.Aborted("etag mismatch")
	}
	f.etag = fmt.Sprintf("etag-%d", f.uploadCalls)
	f.content = req.Content
	return &metadata.UploadResponse{Etag: f.etag}, nil
}

// TestListReceivedSharesBlocksBehindSlowCAS reproduces OCISDEV-1560's
// client-visible symptom one layer down: a reader calling providercache.Get
// (the primitive ListReceivedShares uses per share) blocks for the entire
// duration of a concurrent Add() that is stuck resolving CAS conflicts on
// the same space.
//
// As of the current stable-8.0 (reva pinned 2026-10-14), providercache.Add
// and receivedsharecache's retryPersist no longer hammer storage-system in
// a tight loop -- both route through cas.RetryPersist / cas.NewBackoff
// (github.com/cenkalti/backoff/v5, capped at 50ms) and call
// resyncAfterConflict between attempts instead of retrying blind. That
// closed the "zero backoff, 100 retries back to back" failure mode.
// What it did NOT change: LockSpace/lockUser is still held across every
// attempt, backoff sleep, and resync, so a reader still queues for the
// writer's *entire* CAS resolution, however long that takes -- and per the
// trace (370ffc41d9c9ac...), a single Upload attempt alone can take ~9.7s
// when storage-system itself is slow, dwarfing the 50ms backoff cap. This
// test models that: failed CAS attempts cost a slow network round trip,
// not a retry-count problem, and the reader still blocks for the lot.
func TestListReceivedSharesBlocksBehindSlowCAS(t *testing.T) {
	storageID, spaceID, shareID := "storage1", "space1", "share1"

	const (
		uploadLatency   = 700 * time.Millisecond
		downloadLatency = 300 * time.Millisecond
		roundTrip       = uploadLatency + downloadLatency // one CAS retry's I/O cost ("action")
		failUploads     = 9                                // bootstrap upload + 8 CAS retries before etag matches
	)

	writerHasLock := make(chan struct{})
	storage := &fakeStorage{
		uploadLatency:   uploadLatency,
		downloadLatency: downloadLatency,
		failUploads:     failUploads,
		onFirstDownload: func() { close(writerHasLock) },
	}

	cache := providercache.New(storage, time.Hour)
	ctx := context.Background()

	var writerDur, readerWaitDur time.Duration
	var wg sync.WaitGroup
	wg.Add(2)

	start := time.Now()

	// Writer: Add() takes LockSpace, then loops retrying the CAS upload
	// (providercache.go Add, ~line 163-224) while holding that lock for
	// every retry's network round trip.
	go func() {
		defer wg.Done()
		share := &collaboration.Share{Id: &collaboration.ShareId{OpaqueId: shareID}}
		_ = cache.Add(ctx, storageID, spaceID, shareID, share)
		writerDur = time.Since(start)
	}()

	// Reader: waits until the writer has definitely grabbed LockSpace,
	// then does the same Get() ListReceivedShares does per share
	// (jsoncs3.go:967, providercache.go:295-321). Everything Get() spends
	// before it returns, minus its own work, is pure lock-wait -- this is
	// what the "Grab lock" span measures in production.
	go func() {
		defer wg.Done()
		<-writerHasLock
		waitStart := time.Now()
		_, _ = cache.Get(ctx, storageID, spaceID, shareID, true)
		readerWaitDur = time.Since(waitStart)
	}()

	wg.Wait()

	// Measure the reader's own work in isolation, with no contention left,
	// to show what the "action" actually costs on its own.
	actionStart := time.Now()
	_, _ = cache.Get(ctx, storageID, spaceID, shareID, true)
	readerActionDur := time.Since(actionStart)

	t.Logf("one CAS retry's I/O round trip (upload+download) = %v   <- the 'action'", roundTrip)
	t.Logf("writer held the lock through %d retries           = %v", failUploads, writerDur)
	t.Logf("reader's Get() blocked behind that lock            = %v   <- the 'Grab lock' wait", readerWaitDur)
	t.Logf("reader's Get() alone, no contention                = %v   <- what the action costs by itself", readerActionDur)

	if readerActionDur > roundTrip/2 {
		t.Fatalf("expected the reader's own work to be near-instant, got %v", readerActionDur)
	}
	if readerWaitDur < roundTrip*time.Duration(failUploads)/2 {
		t.Fatalf("expected the reader to be blocked for close to the writer's full retry-storm duration (%v), only blocked for %v", writerDur, readerWaitDur)
	}
}

// TestListReceivedSharesWorkerPoolSelfCollides reproduces the mechanism
// confirmed from the production trace (370ffc41d9c9ac...): ListReceivedShares
// fans out one worker per space (jsoncs3.go:927-940), and several of those
// workers, for the *same* user but *different* spaces, each need to write
// back a discovered share into that user's received-share cache document
// (receivedsharecache.Add). Because receivedsharecache.lockUser locks per
// userID rather than per (userID, spaceID), these writes cannot run in
// parallel even though they touch disjoint spaces -- they queue behind one
// mutex and each pays a full download+CAS-upload round trip in turn.
//
// In the trace this showed up as non-overlapping "Grab lock"+"Add" spans,
// one after another, whose durations summed to ~17.7s of a 48s request.
// This test reproduces that shape directly against the real
// receivedsharecache.Cache: N workers, N distinct spaceIDs, one userID.
func TestListReceivedSharesWorkerPoolSelfCollides(t *testing.T) {
	const (
		userID     = "6d2f6853-bbf4-47e9-952b-4fff9da08770"
		numWorkers = 22 // enough for total wall time to clear 20s at ~1s/worker steady-state
		// downloadLatency only matters for the one-time cold-start sync (the
		// baseline call below); every subsequent per-user Add(), once the
		// cache is warm, costs exactly one Upload on the happy path -- no
		// conflict, no resync -- so uploadLatency alone is the steady-state
		// "action" cost per worker.
		uploadLatency   = 1 * time.Second
		downloadLatency = 300 * time.Millisecond
		roundTrip       = uploadLatency // steady-state per-worker cost once warm
	)

	storage := &fakeStorage{
		uploadLatency:   uploadLatency,
		downloadLatency: downloadLatency,
		failUploads:     0, // no CAS conflicts needed here -- the bug is queueing, not retries
	}
	cache := receivedsharecache.New(storage, time.Hour)
	ctx := context.Background()

	// Baseline: what one Add() costs with nobody else contending.
	baselineStart := time.Now()
	_ = cache.Add(ctx, userID, "space-baseline", &collaboration.ReceivedShare{
		Share: &collaboration.Share{
			Id:         &collaboration.ShareId{OpaqueId: "share-baseline"},
			ResourceId: &provider.ResourceId{StorageId: "storage1", SpaceId: "space-baseline"},
		},
		State: collaboration.ShareState_SHARE_STATE_PENDING,
	})
	baselineDur := time.Since(baselineStart)

	// Simulate numWorkers goroutines of a single ListReceivedShares call,
	// one per space, all discovering a new share for the same user at once.
	durations := make([]time.Duration, numWorkers)
	var wg sync.WaitGroup
	wg.Add(numWorkers)

	start := time.Now()
	for i := 0; i < numWorkers; i++ {
		go func(i int) {
			defer wg.Done()
			spaceID := fmt.Sprintf("space-%d", i)
			shareID := fmt.Sprintf("share-%d", i)
			workerStart := time.Now()
			_ = cache.Add(ctx, userID, spaceID, &collaboration.ReceivedShare{
				Share: &collaboration.Share{
					Id:         &collaboration.ShareId{OpaqueId: shareID},
					ResourceId: &provider.ResourceId{StorageId: "storage1", SpaceId: spaceID},
				},
				State: collaboration.ShareState_SHARE_STATE_PENDING,
			})
			durations[i] = time.Since(workerStart)
		}(i)
	}
	wg.Wait()
	totalWallTime := time.Since(start)

	t.Logf("one Add(), no contention                 = %v   <- the 'action'", baselineDur)
	for i, d := range durations {
		t.Logf("worker %d Add() (different space, same user) = %v", i, d)
	}
	t.Logf("%d workers, wall-clock time for all to finish = %v   <- this is the ListReceivedShares-equivalent total", numWorkers, totalWallTime)

	// Find whichever worker landed 9th in the queue -- goroutine scheduling,
	// not slice index, decides lock acquisition order, so sort by duration
	// rather than assume durations[8] queued 9th. Under pure per-user
	// serialization, the 9th-in-line Add() call should cost ~9 x roundTrip
	// (~9s): 8 other workers' round trips as stall, plus its own ~1s
	// action. This is the exact "9s Grab lock, 1s action, 8s stall" shape
	// reported against the production trace, reproduced here against the
	// unmodified receivedsharecache.Cache.
	sorted := append([]time.Duration(nil), durations...)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i] < sorted[j] })
	const ninthInQueue = 8 // 0-indexed: the 9th-smallest duration
	ninthDur := sorted[ninthInQueue]
	impliedStall := ninthDur - roundTrip
	t.Logf("9th-in-queue worker's Add() call          = %v   <- action %v + stall %v", ninthDur, roundTrip, impliedStall)

	if ninthDur < 8*roundTrip {
		t.Fatalf("expected the 9th-in-queue worker to cost ~9 x %v (~9s: 1s action + 8s stall), only took %v", roundTrip, ninthDur)
	}
	if impliedStall < 7*roundTrip {
		t.Fatalf("expected the 9th-in-queue worker's stall (total minus its own action) to be ~8s, got %v", impliedStall)
	}

	// If the lock were scoped per (userID, spaceID) instead of per userID,
	// these numWorkers writes to disjoint spaces would overlap and the
	// total wall time would stay close to a single round trip. Because the
	// lock is per-user, they serialize: total time should scale with
	// numWorkers, not stay flat -- this is the "ListReceivedShares takes
	// 20+s" claim, reproduced as a measured total, not an estimate.
	const wantMinWallTime = 20 * time.Second
	if totalWallTime < wantMinWallTime {
		t.Fatalf("expected %d per-user Add() calls for different spaces to serialize past %v, completed in %v -- lock contention not reproduced", numWorkers, wantMinWallTime, totalWallTime)
	}
}
