Enhancement: Add new build targets

Make build target `build` used to build a binary twice, the second occurrence having symbols for debugging. We split this step in two and added `build-all` and `build-debug` targets.

- `build-all` now behaves as the previous `build` target, it will generate 2 binaries, one for debug.
- `build-debug` will build a single binary for debugging.

https://github.com/owncloud/ocis/pull/1824
