# Kiteworks Storage Provider for oCIS — Design Spec

**Date:** 2026-05-10
**Author:** Thomas Müller
**Status:** Approved

## Overview

A new oCIS service (`storage-kiteworks`) that exposes a Kiteworks server as a CS3 storage provider. It implements the reva `storage.FS` interface, translating CS3 operations into Kiteworks REST API calls. Each Kiteworks top-level folder becomes a CS3 `StorageSpace`. Authentication is pure OIDC token passthrough — no separate credentials.

## Goals

- Allow oCIS clients (desktop sync, web, mobile) to access files stored on a Kiteworks server.
- Each Kiteworks top-level folder appears as a distinct storage space.
- Owned folders map to `project` spaces; received shares map to `mountpoint` spaces.
- Implement all Kiteworks-supported operations; stub unsupported ones with `errtypes.NotSupported`.

## Non-Goals

- Resumable upload recovery (re-attach after failure).
- NATS event bus emission (deferred to a later iteration).
- Search (`ListFolder` with search query).
- Creating or deleting top-level Kiteworks folders from oCIS (admin-managed lifecycle).

---

## Section 1: Service Structure

### `services/storage-kiteworks/` — oCIS service layer

```
cmd/storage-kiteworks/main.go          # signal context, DefaultConfig, Execute
pkg/command/root.go                    # GetCommands() → Server, Health, Version
pkg/command/server.go                  # wires reva runtime with kiteworks driver
pkg/config/config.go                   # Config struct
pkg/config/defaults.go                 # DefaultConfig()
pkg/revaconfig/config.go               # StorageKiteworksConfigFromStruct()
pkg/server/debug/server.go             # debug HTTP server
pkg/logging/logging.go                 # zerolog wrapper
```

### `vendor/github.com/owncloud/reva/v2/pkg/storage/fs/kiteworks/` — reva storage driver

```
kiteworks.go    # storage.FS implementation; init() registers "kiteworks" driver
client.go       # Kiteworks REST API client (adapted from github.com/owncloud/kwdav pkg/kwlib)
types.go        # Kiteworks data model (FileInfo, DirectoryInfo, UploadResult, etc.)
upload.go       # chunked upload logic
```

The driver is registered via `init()` into reva's `fs/registry` and imported with a blank import in:
```
vendor/github.com/owncloud/reva/v2/pkg/storage/fs/loader/loader.go
```

The service binary is registered in the main oCIS multi-service runtime (same pattern as `storage-users`, `storage-system`). It is opt-in: the service only starts if `STORAGE_KITEWORKS_ENDPOINT` is configured, so existing deployments are unaffected.

---

## Section 2: CS3 → Kiteworks Operation Mapping

### Storage Spaces

| CS3 operation | Kiteworks REST | Notes |
|---|---|---|
| `ListStorageSpaces` | `GET /rest/folders/top` | Each folder → one `StorageSpace`. See space type mapping below. |
| `CreateStorageSpace` | — | `errtypes.NotSupported` |
| `UpdateStorageSpace` | — | `errtypes.NotSupported` |
| `DeleteStorageSpace` | — | `errtypes.NotSupported` |
| `CreateHome` | — | `errtypes.NotSupported` (deprecated) |
| `GetHome` | — | `errtypes.NotSupported` (deprecated) |

**Space type mapping** (from `FileInfo.IsSharedWithUser()`):

- `IsSharedWithUser() == true` (i.e. `IsShared == true` and `ParentID == nil || "0"`) → `SpaceType = "mountpoint"` (received share)
- `IsSharedWithUser() == false` → `SpaceType = "project"` (owned folder)

The `StorageSpace.Id.OpaqueId` is set to the Kiteworks folder ID. No extra API calls are needed to determine ownership — `GET /rest/folders/top` provides all required fields.

### File & Folder CRUD

| CS3 operation | Kiteworks REST | Notes |
|---|---|---|
| `GetMD` (by ID) | `GET /rest/folders/{id}` or `GET /rest/files/{id}` | Try folders endpoint first; on 404 fall back to files endpoint. CS3 `ResourceId` does not carry a type discriminator. |
| `GetMD` (by path) | `GET /rest/query` (search by path) | Fall back to path-based search |
| `ListFolder` | `GET /rest/folders/{id}/children` | |
| `CreateDir` | `POST /rest/folders/{id}/folders` | |
| `Delete` | `DELETE /rest/folders/{id}` or `DELETE /rest/files/{id}` | |
| `Move` | `POST /rest/files/actions/move` | Works for files and folders |
| `TouchFile` | Initiate upload + single empty chunk | See upload section |
| `GetPathByID` | `GET /rest/files/{id}` or `/rest/folders/{id}` | Return `path` field |
| `CreateReference` | — | `errtypes.NotSupported` |

### Upload & Download

| CS3 operation | Kiteworks REST | Notes |
|---|---|---|
| `InitiateUpload` | `POST /rest/folders/{id}/actions/initiateUpload` | Returns `uploadID` + `uploadURI` in map |
| `Upload` | `UploadChunk` per chunk | Chunks of `ChunkSize` bytes (default 5 MB); last chunk sets `isLastChunk=true` |
| `Download` | `GET /rest/files/{id}/content` | `Range` header passed through for partial reads |

TUS (`ComposableFS`) is not implemented — the Kiteworks chunked protocol is not TUS-compatible. oCIS falls back to the direct `Upload` path.

### Revisions

Stubbed with `errtypes.NotSupported` initially. Promote to real implementation if the Kiteworks REST API exposes version history (to be confirmed against the OpenAPI spec at implementation time).

### Recycle Bin

Stubbed with `errtypes.NotSupported` initially. Promote to real implementation if the Kiteworks REST API exposes trash operations.

### Grants

| CS3 operation | Kiteworks REST | Notes |
|---|---|---|
| `ListGrants` | File/folder metadata `Permission` array | Map Kiteworks `Permission` → CS3 `Grant` |
| `AddGrant` | Kiteworks permissions endpoint | Map CS3 role → Kiteworks permission ID |
| `UpdateGrant` | Kiteworks permissions endpoint | |
| `RemoveGrant` | Kiteworks permissions endpoint | |
| `DenyGrant` | — | `errtypes.NotSupported` (no explicit deny ACL in Kiteworks) |

### Metadata, Quota & Locks

| CS3 operation | Kiteworks REST | Notes |
|---|---|---|
| `GetQuota` | `GET /rest/users/me` | Returns `QuotaInfo.Allowed` / `QuotaInfo.Used` |
| `SetArbitraryMetadata` | — | `errtypes.NotSupported` |
| `UnsetArbitraryMetadata` | — | `errtypes.NotSupported` |
| `GetLock` | — | `errtypes.NotSupported` |
| `SetLock` | — | `errtypes.NotSupported` |
| `RefreshLock` | — | `errtypes.NotSupported` |
| `Unlock` | — | `errtypes.NotSupported` |

---

## Section 3: Configuration & Authentication

### Config struct (`pkg/config/config.go`)

```go
type KiteworksDriver struct {
    Endpoint  string `yaml:"endpoint"   env:"STORAGE_KITEWORKS_ENDPOINT"   desc:"Base URL of the Kiteworks server, e.g. https://kiteworks.example.com"`
    Insecure  bool   `yaml:"insecure"   env:"STORAGE_KITEWORKS_INSECURE"    desc:"Skip TLS certificate verification."`
    MountID   string `yaml:"mount_id"   env:"STORAGE_KITEWORKS_MOUNT_ID"    desc:"Mount ID of this storage provider."`
    ChunkSize int64  `yaml:"chunk_size" env:"STORAGE_KITEWORKS_CHUNK_SIZE"  desc:"Upload chunk size in bytes. Default 5242880 (5 MB)."`
}
```

No credentials in config. Authentication is token passthrough only.

### Authentication flow

1. Every CS3 gRPC request carries the user's OIDC token in the context (extracted by reva's auth middleware).
2. The driver extracts it via `ctxpkg.ContextGetToken(ctx)`.
3. Each outbound Kiteworks REST request includes `Authorization: Bearer <token>`.
4. An `APIClient` is constructed per-request with the token — no shared client state across users.

**Deployment requirement:** The Kiteworks server must accept the OIDC tokens issued by oCIS's IDP (or they share a common token issuer). This is the expected topology when oCIS delegates identity to Kiteworks' own IDP.

### Reva driver config mapping (`pkg/revaconfig/config.go`)

```go
func Kiteworks(cfg *config.Config) map[string]interface{} {
    return map[string]interface{}{
        "endpoint":   cfg.Driver.Endpoint,
        "insecure":   cfg.Driver.Insecure,
        "chunk_size": cfg.Driver.ChunkSize,
    }
}
```

---

## Section 4: Upload Flow

### `InitiateUpload`

Calls `POST /rest/folders/{parentID}/actions/initiateUpload` with:
- `filename` from `metadata["filename"]`
- `totalSize = uploadLength`
- `numberOfChunks = ceil(uploadLength / ChunkSize)` (minimum 1)

Returns a map:
```go
map[string]string{
    "uploadID":  result.ID,
    "uploadURI": result.URI,
}
```

### `Upload`

1. Reads `req.Body` in `ChunkSize` slices.
2. For each slice, calls `UploadChunk(uploadURI, filename, body, chunkIndex, chunkData, isLastChunk)`.
3. The response from the final chunk is the committed `FileInfo`, which is converted to `*provider.ResourceInfo`.

### `TouchFile`

Equivalent to uploading a zero-byte file:
1. `InitiateUpload` with `totalSize=0, numberOfChunks=1`.
2. Single `UploadChunk` call with an empty body and `isLastChunk=true`.

### Chunk size default

`5 * 1024 * 1024` bytes (5 MB). Configurable via `STORAGE_KITEWORKS_CHUNK_SIZE`.

---

## Section 5: Integration & Testing

### Integration into oCIS binary

- Add blank import in `vendor/.../reva/v2/pkg/storage/fs/loader/loader.go`.
- Register the service in the oCIS multi-service runtime (same as `storage-users`).
- Service is opt-in: starts only when `STORAGE_KITEWORKS_ENDPOINT` is set.

### Testing

- **Unit tests** for the client layer using `httptest.NewServer` mock (pattern: `nextcloud_server_mock.go`).
- **Unit tests** for each `storage.FS` method verifying correct REST endpoint is called and correct CS3 response is returned.
- **Ginkgo v2 test suite** bootstrap: `TestKiteworks(t *testing.T)` with `RegisterFailHandler(Fail)` / `RunSpecs`.
- No acceptance tests in v1 (require live Kiteworks instance).

### Out of scope for v1

- Resumable upload recovery (re-attach to an existing Kiteworks upload session after failure)
- NATS event bus emission
- Search support in `ListFolder`
- Versioning and recycle bin (stubbed; implement if API supports them)
