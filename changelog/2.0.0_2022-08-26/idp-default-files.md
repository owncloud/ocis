Enhancement: Generate signing key and encryption secret

The idp service now automatically generates a signing key and encryption secret when they don't exist.
This will enable service restarts without invalidating existing sessions.

https://github.com/owncloud/ocis/issues/3909
https://github.com/owncloud/ocis/pull/4022
