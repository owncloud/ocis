Enhancement: Update OCIS Runtime

- enhances the overall behavior of our runtime
- runtime `db` file configurable
- two new env variables to deal with the runtime
- `RUNTIME_DB_FILE` and `RUNTIME_KEEP_ALIVE`
- `RUNTIME_KEEP_ALIVE` defaults to `false` to provide backwards compatibility
- if `RUNTIME_KEEP_ALIVE` is set to `true`, if a supervised process terminates the runtime will attempt to start with the same environment provided.

https://github.com/owncloud/ocis/pull/1108
