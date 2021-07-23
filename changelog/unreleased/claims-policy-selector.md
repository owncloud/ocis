Enhancement: Proxy: Add claims policy selector

Using the proxy config file, it is now possible to let let the IdP determine the routing policy by sending an `ocis.routing.policy` claim. Its value will be used to determine the set of routes for the logged in user.

https://github.com/owncloud/ocis/pull/2248
