Bugfix: Never expire the storage-users ID cache

The storage-users ID cache holds the authoritative id<->path index used by the
storage provider. It was given a cache TTL (24 minutes by default, and settable
via `OCIS_CACHE_TTL` / `STORAGE_USERS_ID_CACHE_TTL`), which the cache layer
applies as a per-write expiry for in-memory stores and as the bucket-wide MaxAge
for the nats-js-kv store. Once entries aged out, the provider could no longer
resolve existing nodes: with the POSIX driver files and folders appeared to
vanish (and sync clients were told to delete them locally), and with the
decomposed drivers they were repeatedly re-resolved from disk.

Because expiring this index is a data-availability problem rather than a tuning
preference, the TTL is no longer configurable for the ID cache: the
`STORAGE_USERS_ID_CACHE_TTL` setting is removed and the reva cache TTL for the
ID cache is pinned to 0 (no expiry), regardless of `OCIS_CACHE_TTL`. The
file-metadata cache (a real cache) keeps its TTL.

Note: existing `nats-js-kv` ID-cache buckets created with a previous non-zero
MaxAge keep it until the bucket is recreated.

https://github.com/owncloud/ocis/pull/12416
