Bugfix: Fix 403 in docs pipeline

Docs pipeline was not routed through our proxies which could lead to requests being blacklisted

https://github.com/owncloud/ocis/issues/7509
https://github.com/owncloud/ocis/pull/7511
