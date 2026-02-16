Bugfix: Fix `OCIS_RUN_SERVICES`

`OCIS_RUN_SERVICES` was introduced as successor to `OCIS_RUN_EXTENSIONS` because
we wanted to call oCIS "core" extensions services. We kept `OCIS_RUN_EXTENSIONS` for
backwards compatibility reasons.

It turned out, that setting `OCIS_RUN_SERVICES` has no effect since introduced. `OCIS_RUN_EXTENSIONS`.
`OCIS_RUN_EXTENSIONS` was working fine all the time.

We now fixed `OCIS_RUN_SERVICES`, so that you can use it as a equivalent replacement for `OCIS_RUN_EXTENSIONS`

https://github.com/owncloud/ocis/pull/4133
