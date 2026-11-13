Bugfix: Fix slow builds with GNU Make 4.4

The file `.make/recursion.mk` now expands the `GO`, `GOPATH`, `GOHOSTOS`,
`GOHOSTARCH` and `GOHOSTARM` variables one time only, before its `export`.
Before, each use of these variables ran a shell command again. With GNU Make
4.4, these commands ran again and again. A `make` call took about two minutes
before it started the build, and it printed thousands of "env: 'GOHOSTOS': No
such file or directory" errors. Ubuntu 26.04 and the `ubuntu-latest` GitHub
runner have GNU Make 4.4.

https://github.com/owncloud/ocis/pull/13458
