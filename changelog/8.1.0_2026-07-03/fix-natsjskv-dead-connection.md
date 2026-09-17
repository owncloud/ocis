Bugfix: Recover from permanently-closed NATS connections in the nats-js-kv store

The `nats-js-kv` go-micro store plugin's `hasConn()` only checked whether the
connection object was non-nil, not whether it was still alive. Once the
underlying NATS client exhausted its reconnect attempts (e.g. a NATS pod
restart longer than the client's reconnect window), the connection stayed
non-nil but permanently closed. Because connection initialization is gated on
`!hasConn()`, it never re-ran, so every subsequent KV operation failed with
`nats: connection closed` until the affected pod was restarted.

This surfaced as several user-visible failures backed by the NATS KV cache,
e.g. all spaces becoming invisible (`storage-users` `ListStorageSpaces`) and
download failures from missing signing keys (`ocs`).

The store plugin now treats a closed connection as no connection, so the next
operation transparently re-initializes it.

https://github.com/owncloud/ocis/pull/12402
