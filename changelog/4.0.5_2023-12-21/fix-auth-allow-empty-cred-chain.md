Bugfix: fix reva config of frontend service to avoid misleading error logs

We set an empty Credentials chain for the frontend service now. In ocis all
non-reva token authentication is handled by the proxy. This avoids irritating
error messages about the missing 'auth-bearer' service.

https://github.com/owncloud/ocis/pull/7934
https://github.com/owncloud/ocis/pull/7453
https://github.com/cs3org/reva/pull/4396
https://github.com/cs3org/reva/pull/4241
https://github.com/owncloud/ocis/issues/6692
