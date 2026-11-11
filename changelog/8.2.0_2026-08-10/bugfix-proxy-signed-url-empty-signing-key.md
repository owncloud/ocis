Bugfix: Prevent panic in presigned URL auth when signing key is missing

The presigned-URL authenticator read the per-user signing key from a store and
indexed the first record without checking that the result was non-empty. When
the store returns an empty result without an error — which the `noop` store
does, and `OCIS_CACHE_STORE=noop` also switches the presigned-URL signing-key
store to `noop` — every signed-URL request panicked with `index out of range
[0] with length 0` (recovered per request as an HTTP 502), so all signed-URL
downloads failed. The authenticator now guards the lookup and returns a normal
authentication error instead of panicking; it also no longer returns a nil
error when the stored key is empty.

https://github.com/owncloud/ocis/issues/12384
https://github.com/owncloud/ocis/pull/12418
