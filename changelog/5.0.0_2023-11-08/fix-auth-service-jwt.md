Bugfix: Fix the auth service env variable

We the auth service env variable to the service specific name. Before it was configurable via `AUTH_MACHINE_JWT_SECRET` and now is configurable via `AUTH_SERVICE_JWT_SECRET`.

https://github.com/owncloud/ocis/pull/7523
