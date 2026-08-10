Bugfix: Re-verify access tokens with an expired userinfo cache entry

The proxy authenticated requests against a cached userinfo entry. When the
cached entry's expiry had passed, the request was rejected outright instead of
re-verifying the access token. The cached expiry is not always the real token
expiry — for opaque access tokens it falls back to the configured cache TTL —
so a still-valid token could be rejected just because its cache entry aged out,
causing spurious logouts (most noticeable under frequent PROPFIND requests). The
proxy now treats an expired cache entry as a cache miss and re-verifies the
access token, while a genuinely expired or invalid token is still rejected by
the re-verification.

https://github.com/owncloud/ocis/pull/12414
