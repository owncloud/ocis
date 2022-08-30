Enhancement: Refactor the proxy service

The routes of the proxy service now have a "unprotected" flag. This is used by the authentication middleware to determine if the request needs to be blocked when missing authentication or not. 

https://github.com/owncloud/ocis/issues/4401
https://github.com/owncloud/ocis/pull/4461
