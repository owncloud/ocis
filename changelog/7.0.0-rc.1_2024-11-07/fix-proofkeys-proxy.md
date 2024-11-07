Bugfix: Respect proxy url when validating proofkeys

We fixed a bug where the proxied wopi URL was not used when validating proofkeys. This caused the validation to fail when the proxy was used.

https://github.com/owncloud/ocis/pull/1234
