Security: Build with a newer Go version

We've raised the Go version to 1.27.3. This version fixes security issues in
the Go standard library. We've also upgraded golangci-lint to v2.15.0, because
the old version does not work with Go 1.27.

https://github.com/owncloud/ocis/pull/13592
