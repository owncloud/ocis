Bugfix: Do not expire the storage-users ID cache by default

The storage-users ID cache holds the authoritative id<->path index used by the
storage provider. Its default TTL was 24 minutes, which the cache layer applies
as a per-write expiry for in-memory stores and as the bucket-wide MaxAge for the
nats-js-kv store. Once entries aged out, the provider could no longer resolve
existing nodes: with the POSIX driver files and folders appeared to vanish, and
with the decomposed drivers they were repeatedly re-resolved from disk. The ID
cache is an index rather than transient data, so the default TTL is removed (0 =
no expiry). The `STORAGE_USERS_ID_CACHE_TTL` setting is kept for operators who
explicitly want one.

Note: existing `nats-js-kv` ID-cache buckets created with the previous 24m MaxAge
keep that MaxAge until the bucket is recreated; the fix prevents new deployments
from getting the harmful default.

https://github.com/owncloud/ocis/pull/12416
