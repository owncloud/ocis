Security: Bump go-jose and golang.org/x/net dependencies

We've bumped go-jose/v4 to v4.1.6 to fix CVE-2026-116587, a denial of service
when parsing JWKs that contain RSA keys, and golang.org/x/net to v0.62.0 to
fix GO-2026-6860 and GO-2026-6866 in its html package.

https://github.com/owncloud/ocis/pull/13256
