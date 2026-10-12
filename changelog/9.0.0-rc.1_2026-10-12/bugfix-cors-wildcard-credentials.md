Bugfix: Do not allow credentials for wildcard CORS origins

CORS responses no longer send `Access-Control-Allow-Credentials: true` when the
configured origins permit any origin: an empty list, `*`, or a wildcard pattern
such as `https://*`. In those cases credentials are dropped and a warning is
logged. This applies across all CORS paths (the shared HTTP middleware, the reva
interceptor used by the frontend and ocm services, ocdav, and the tus upload
handler). Explicitly configured origins are unaffected and keep credentials, so
intended cross-origin setups continue to work.

https://github.com/owncloud/ocis/pull/12983
