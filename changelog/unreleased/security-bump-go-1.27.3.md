Security: Build with Go 1.27.3

We've raised the Go version to 1.27.3. This version fixes CVE-2026-78667,
CVE-2026-78669 and CVE-2026-111943 in the Go standard library. We've also
upgraded golangci-lint to v2.15.0, because the old version cannot read the
export data of Go 1.27.

https://github.com/owncloud/ocis/pull/13592
