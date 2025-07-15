Enhancement: Rewrite of the request authentication middleware

There were some flaws in the authentication middleware which were resolved by this rewrite.
This rewrite also introduced the need to manually mark certain paths as "unprotected" if
requests to these paths must not be authenticated.

https://github.com/owncloud/ocis/pull/4374
