Enhancement: Refactor graph API

We refactored the `/graph/v1.0/` endpoint which now relies on the internal acces token fer authentication, getting rid of any LDAP or OIDC code to authenticate requests. This allows using the graph api when using basic auth or any other auth mechanism provided by the CS3 auth providers / reva gateway / ocis proxy.

https://github.com/owncloud/ocis/pull/2277
