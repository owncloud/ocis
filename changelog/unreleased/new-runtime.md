Change: Switch over to a new custom-built runtime

We moved away from using the go-micro runtime and are now using [our own runtime](https://github.com/refs/pman).
This allows us to spawn service processes even when they are using different versions of go-micro. On top of that we
now have the commands `ocis list`, `ocis kill` and `ocis run` available for service runtime management.

https://github.com/owncloud/ocis/pull/287
