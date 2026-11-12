Security: Bump idp frontend dependencies

We've refreshed the idp frontend lockfile to pick up fixes for two transitive
dependencies. `follow-redirects` moves to 1.16.2, which fixes CVE-2026-111755,
CVE-2026-111756 and CVE-2026-116689, three information disclosures via
cross-origin redirects. `decode-uri-component` moves to 0.5.1 through
`query-string` 9.5.1, which fixes CVE-2026-118664, a denial of service via
malformed percent-encoded input.

Both fixes were already allowed by the version ranges the idp declares, so no
dependency was bumped across a major version.

https://github.com/owncloud/ocis/pull/13345
