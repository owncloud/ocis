Change: Move ocis default config to root level

Tags: ocis

We moved the tracing config to the `root` flagset so that they are parsed on all commands. We also introduced a `JWTSecret` flag in the root flagset, in order to apply a common default JWTSecret to all services that have one.

https://github.com/owncloud/ocis/pull/842
https://github.com/owncloud/ocis/pull/843
