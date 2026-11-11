Enhancement: Add TLS support for the frontend stat cache store connection

The frontend service was the only remaining place where a `nats-js-kv` store
connection could not be secured. Its OCS stat cache forwarded the store type,
nodes, database, table, TTL and credentials to reva, but not the TLS settings,
so the connection stayed plaintext with no operator toggle to change it. The
stat cache now honours `OCIS_CACHE_ENABLE_TLS`, `OCIS_CACHE_TLS_INSECURE` and
`OCIS_CACHE_TLS_ROOT_CA_CERTIFICATE` like every other cache and store.

https://github.com/owncloud/ocis/pull/12931
