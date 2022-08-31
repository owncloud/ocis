Bugfix: remove legacy accounts proxy routes

We've removed the legacy accounts routes from the proxy default config.
There were no longer used since the switch to IDM as the default user
backend. Also accounts is no longer part of the oCIS binary and therefore
should not be part of the proxy default route config.

https://github.com/owncloud/ocis/pull/3831
