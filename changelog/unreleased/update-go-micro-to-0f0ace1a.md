Enhancement: Update go-micro to v3.5.1-0.20210217182006-0f0ace1a44a9

- We updated from go micro v2 (v2.9.1) go-micro v3 (v3.5.1 edge).
- oCIS runtime is now aware of `MICRO_LOG_LEVEL` and is set to `error` by default. This decision was made because ownCloud, as framework builders, want to log everything oCIS related and hide everything unrelated by default. It can be re-enabled by setting it to a log level other than `error`. i.e: `MICRO_LOG_LEVEL=info`.
- Updated `protoc-gen-micro` to the [latest version](https://github.com/asim/go-micro/tree/master/cmd/protoc-gen-micro).
- We're using Prometheus wrappers from go-micro.

https://github.com/owncloud/ocis/pull/1670
https://github.com/asim/go-micro/pull/2126
