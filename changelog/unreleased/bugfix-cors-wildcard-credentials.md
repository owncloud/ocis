Bugfix: Do not allow credentials for wildcard CORS origins

The shared CORS middleware no longer sends `Access-Control-Allow-Credentials: true`
when origins are unrestricted (an empty list or `*`); credentials are dropped and
a warning is logged. Explicitly configured origins are unaffected and keep
credentials, so intended cross-origin setups continue to work.

https://github.com/owncloud/ocis/pull/
