Bugfix: Return metadata client init errors from the settings store

The settings metadata store initialized its metadata client lazily. When
initialization failed, for example because the settings metadata space was
owned by a different user than the configured system user, `Store.Init()`
logged the error but left the client nil and returned normally. Every
subsequent store call then dereferenced the nil client and panicked with a nil
pointer dereference, which masked the real cause and surfaced as opaque proxy
500s during login.

`Store.Init()` now returns the underlying initialization error and every public
store method short-circuits on it, so callers receive a specific, actionable
error instead of a panic. The client is left nil on failure so the next call
retries initialization.

https://github.com/owncloud/ocis/pull/12680
