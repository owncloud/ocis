Bugfix: forward basic auth to OpenID connect token authentication endpoint

When using `PROXY_ENABLE_BASIC_AUTH=true` we now forward request to the idp instead of trying to authenticate the request ourself.

https://github.com/owncloud/ocis/issues/2095
https://github.com/owncloud/ocis/issues/2094