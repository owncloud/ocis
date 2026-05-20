# Design: Eliminate Redundant Stat on WebDAV GET

**Date:** 2026-05-20
**Repo:** `github.com/owncloud/reva`
**Tracking:** `download-perf-investigation.md` (root of ocis repo)

---

## Problem

Every WebDAV GET request through `remote.php/webdav/{userid}` pays an extra gRPC `Stat` round-trip before the download begins. This adds sequential latency that upload does not pay, making downloads measurably slower than other operations.

Root cause: `handleGet` in `ocdav/get.go` calls `client.Stat()` at line 72 solely to:
1. Detect directories (return 200 empty body)
2. Detect files still being processed (return 425)

Then immediately calls `client.InitiateFileDownload()` — which internally resolves the same node again. The node is loaded and permissions checked twice; `AsResourceInfo` (with full xattr reads including checksums, grants, share-types) is called twice.

---

## Solution

Move the directory and processing guards into `InitiateFileDownload` at the storage provider layer. `handleGet` drops its `Stat` call entirely.

This collapses three sequential gRPC calls into two for every download:

| Before | After |
|--------|-------|
| `ListStorageSpaces` | `ListStorageSpaces` |
| `Stat` (full AsResourceInfo) | ~~`Stat`~~ |
| `InitiateFileDownload` | `InitiateFileDownload` (with guards) |

---

## Design

### 1. New `FS` interface method — `pkg/storage/storage.go`

```go
// StatForDownload returns the resource type and processing state for ref.
// Allows InitiateFileDownload to gate downloads without a full Stat.
StatForDownload(ctx context.Context, ref *provider.Reference) (provider.ResourceType, bool, error)
```

**Decomposedfs implementation** (`pkg/storage/utils/decomposedfs/decomposedfs.go`):

```go
func (fs *Decomposedfs) StatForDownload(ctx context.Context, ref *provider.Reference) (provider.ResourceType, bool, error) {
    n, err := fs.lu.NodeFromResource(ctx, ref)
    if err != nil {
        return provider.ResourceType_RESOURCE_TYPE_INVALID, false, err
    }
    if !n.Exists {
        return provider.ResourceType_RESOURCE_TYPE_INVALID, false, errtypes.NotFound(n.ID)
    }
    return n.Type(ctx), n.IsProcessing(ctx), nil
}
```

- `n.Type(ctx)` — cheap in-memory read from node struct after load
- `n.IsProcessing(ctx)` — one `getxattr` syscall for `prefixes.StatusPrefix`
- No `AsResourceInfo`, no etag, no checksum reads, no grant enumeration

**All other FS backends** get a minimal stub returning the safe default (allow download to proceed; data transfer fails naturally if ref is invalid):

```go
func (fs *XxxFs) StatForDownload(ctx context.Context, ref *provider.Reference) (provider.ResourceType, bool, error) {
    return provider.ResourceType_RESOURCE_TYPE_FILE, false, nil
}
```

Backends requiring stubs (identified by `GetMD` presence as the FS interface marker):
- `pkg/storage/fs/owncloudsql/owncloudsql.go` — `owncloudsqlfs`
- `pkg/storage/fs/nextcloud/nextcloud.go` — `StorageDriver`
- `pkg/storage/fs/cephfs/cephfs.go` — `cephfs`
- `pkg/storage/fs/hello/hello.go` — `hellofs`
- `pkg/storage/fs/s3/s3.go` — `s3FS`
- `pkg/storage/utils/eosfs/eosfs.go` — `eosfs`
- `pkg/storage/utils/localfs/localfs.go` — `localfs`
- `pkg/storage/utils/middleware/middleware.go` — middleware wrapper (must delegate to inner FS)

### 2. `storageprovider.go` — `InitiateFileDownload`

**File:** `internal/grpc/services/storageprovider/storageprovider.go`

Insert before the URL construction:

```go
resourceType, isProcessing, err := s.storage.StatForDownload(ctx, req.Ref)
if err != nil {
    return &provider.InitiateFileDownloadResponse{
        Status: status.NewInternal(ctx, err.Error()),
    }, nil
}
if resourceType == provider.ResourceType_RESOURCE_TYPE_CONTAINER {
    return &provider.InitiateFileDownloadResponse{
        Status: status.NewInvalidArg(ctx, "resource is a directory"),
    }, nil
}
if isProcessing {
    return &provider.InitiateFileDownloadResponse{
        Status: status.NewOK(ctx),
        Opaque: utils.AppendPlainToOpaque(nil, "status", "processing"),
    }, nil
}
```

No proto changes needed — `Opaque` is already present on `InitiateFileDownloadResponse` and `status=processing` is an existing convention.

### 3. `handleGet` in `ocdav/get.go`

**Remove** (lines 69–91):
- `StatRequest` construction
- `client.Stat()` call and error handling
- Directory type check → 200 empty body
- Processing status check → 425

**Add** after `InitiateFileDownload` succeeds:

```go
// directory guard (was previously caught by Stat)
if dRes.Status.Code == rpc.Code_CODE_INVALID_ARGUMENT {
    w.Header().Set("Content-Length", "0")
    w.WriteHeader(http.StatusMethodNotAllowed)
    return
}

// processing guard (was previously caught by Stat)
if status := utils.ReadPlainFromOpaque(dRes.GetOpaque(), "status"); status == "processing" {
    w.WriteHeader(http.StatusTooEarly)
    return
}
```

Note: the directory case changes from `200 empty body` to `405 Method Not Allowed`. This is a minor behaviour improvement — a GET on a directory is not a valid WebDAV operation.

---

## Scope

All changes are in `github.com/owncloud/reva`. No proto changes. No ocis-level changes except updating `go.mod` to point at the new reva commit once the PR is merged.

Files touched in reva:
| File | Change |
|------|--------|
| `pkg/storage/storage.go` | Add `StatForDownload` to `FS` interface |
| `pkg/storage/utils/decomposedfs/decomposedfs.go` | Implement `StatForDownload` |
| `pkg/storage/fs/owncloudsql/owncloudsql.go` | Stub |
| `pkg/storage/fs/nextcloud/nextcloud.go` | Stub |
| `pkg/storage/fs/cephfs/cephfs.go` | Stub |
| `pkg/storage/fs/hello/hello.go` | Stub |
| `pkg/storage/fs/s3/s3.go` | Stub |
| `pkg/storage/utils/eosfs/eosfs.go` | Stub |
| `pkg/storage/utils/localfs/localfs.go` | Stub |
| `pkg/storage/utils/middleware/middleware.go` | Delegate to inner FS |
| `internal/grpc/services/storageprovider/storageprovider.go` | Call `StatForDownload` in `InitiateFileDownload` |
| `internal/http/services/owncloud/ocdav/get.go` | Remove `Stat` call, handle new error codes |

---

## Testing

- Existing ocdav GET handler tests cover directory and processing cases — update expected error source (from `Stat` mock to `InitiateFileDownload` mock)
- Add unit test for `Decomposedfs.StatForDownload` covering: file, directory, processing file, not-found
- Add unit test for `storageprovider.InitiateFileDownload` covering: directory returns `CODE_INVALID_ARGUMENT`, processing returns opaque `status=processing`, normal file returns protocols

---

## Expected Impact

For a 1MB download over a fast local network, the current latency breakdown is dominated by sequential gRPC hops, not transfer time. Removing `Stat` eliminates:
- One full gRPC round-trip (gateway → storage provider → node resolution → permission check → `AsResourceInfo`)
- Multiple `getxattr` syscalls (etag, SHA1, MD5, Adler32 checksums, grants enumeration)
- One redundant `NodeFromResource` call and one redundant `AssemblePermissions` call

Uploads are already at the lower bound. This change brings download latency to the same baseline.
