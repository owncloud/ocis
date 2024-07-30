Bugfix: Fix crash on empty tracing provider

We have fixed a bug that causes a crash when OCIS_TRACING_ENABLED is set to true, but no
tracing Endpoints or Collectors have been provided.a

https://github.com/owncloud/ocis/pull/9622
https://github.com/owncloud/ocis/issues/7012
