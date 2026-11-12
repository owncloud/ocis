Security: Bump vulnerable dependencies

Bumped `golang.org/x/net` to v0.62.0, `github.com/go-jose/go-jose/v4` to v4.1.6
and `github.com/klauspost/compress` to v1.20.2. The last one fixes GO-2026-6865,
a panic on two-block frames.

https://github.com/owncloud/ocis/pull/13300
https://github.com/owncloud/ocis/pull/13365
