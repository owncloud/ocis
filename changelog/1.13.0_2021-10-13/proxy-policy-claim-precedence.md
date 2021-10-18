Enhancement: allow overriding the cookie based route by claim

When determining the routing policy we now let the claim override the cookie so that users are routed to the correct backend after login.

https://github.com/owncloud/ocis/pull/2508