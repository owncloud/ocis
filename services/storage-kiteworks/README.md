# Storage-Kiteworks

The `storage-kiteworks` service exposes a [Kiteworks](https://kiteworks.com/) server as a CS3 storage provider within oCIS.

Each Kiteworks top-level folder becomes a distinct CS3 `StorageSpace`:

- Folders owned by the authenticated user → space type `"project"`
- Folders received as shares from other users → space type `"mountpoint"`

Authentication uses pure OIDC token passthrough — no separate Kiteworks credentials are required.

## Service Startup

The service is **opt-in** and does not start automatically. To enable it, add it to `OCIS_ADD_RUN_SERVICES`:

```bash
OCIS_ADD_RUN_SERVICES=storage-kiteworks \
STORAGE_KITEWORKS_ENDPOINT=https://kiteworks.example.com \
STORAGE_KITEWORKS_MOUNT_ID=<uuid> \
./ocis server
```

Existing deployments are unaffected when `storage-kiteworks` is omitted from `OCIS_ADD_RUN_SERVICES`.

## Configuration

### Required Environment Variables

| Variable | Description |
|---|---|
| `STORAGE_KITEWORKS_ENDPOINT` | Base URL of the Kiteworks server, e.g. `https://kiteworks.example.com` |
| `STORAGE_KITEWORKS_MOUNT_ID` | Mount ID for this storage provider (must be a stable UUID) |

### Optional Environment Variables

| Variable | Default | Description |
|---|---|---|
| `STORAGE_KITEWORKS_CHUNK_SIZE` | `5242880` | Upload chunk size in bytes (default 5 MB) |
| `STORAGE_KITEWORKS_INSECURE` | `false` | Skip TLS certificate verification. Development only — never use in production |
| `STORAGE_KITEWORKS_GRPC_ADDR` | `127.0.0.1:9285` | Bind address of the gRPC service |
| `STORAGE_KITEWORKS_DEBUG_ADDR` | `127.0.0.1:9289` | Bind address of the debug/metrics server |

For the full list of environment variables see the generated documentation.

## Ports

The service uses port range **9285-9289**:

| Port | Usage |
|---|---|
| 9285 | gRPC storage provider endpoint |
| 9289 | Debug server (metrics, health, pprof, zpages) |

## Health

The service responds to `ocis storage-kiteworks health` and exposes `/healthz` and `/readyz` on the debug port.

## Metrics

Prometheus metrics are exposed at `http://<STORAGE_KITEWORKS_DEBUG_ADDR>/metrics` under the `ocis_storage_kiteworks` namespace.

## Dependency

The Kiteworks storage driver lives in [owncloud/reva](https://github.com/owncloud/reva) at `pkg/storage/fs/kiteworks/`. This service requires a reva build that includes that driver.
