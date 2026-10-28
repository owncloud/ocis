Security: Bump go-jose dependency

We've bumped the go-jose/v4 dependency to v4.1.6 to fix CVE-2026-116587, a
denial of service when parsing JWKs that contain RSA keys.

https://github.com/owncloud/ocis/pull/PLACEHOLDER
