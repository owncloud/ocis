Bugfix: Fix authentication for autoprovisioned users

We've fixed an issue in the proxy, which made the first http request of an
autoprovisioned user fail.

https://github.com/owncloud/ocis/issues/4616
