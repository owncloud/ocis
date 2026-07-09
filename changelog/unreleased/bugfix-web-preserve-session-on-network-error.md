Bugfix: Keep the web session on transient network errors during context init

We've fixed a bug in the web frontend where a temporary connectivity failure (e.g. a
fast page reload cancelling in-flight requests, or a short offline period) while
re-establishing the authenticated context on page load was treated like an expired or
invalid token. This forced a full logout and redirected the user to the "trouble
connecting to the login service" page. Network/connectivity errors now preserve the
existing session, while genuine authentication rejections (e.g. a 401) still trigger
the regular logout flow as before.

https://github.com/owncloud/web/issues/12980
