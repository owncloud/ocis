Enhancement: Add TLS support for NATS store and registry connections

All `nats-js-kv` store, cache, and service registry connections now support
TLS. Configure via `OCIS_CACHE_ENABLE_TLS`, `OCIS_PERSISTENT_STORE_ENABLE_TLS`,
and `MICRO_REGISTRY_ENABLE_TLS`, with corresponding `*_TLS_INSECURE` and
`*_TLS_ROOT_CA_CERTIFICATE` variants per connection type.

https://github.com/owncloud/ocis/pull/12765
