Bugfix: CSP rules for silent token refresh in iframe

When renewing the access token silently web needs to be opened in an iframe. This was previously blocked by a restrictive iframe CSP rule in the `Secure` middleware and has now been fixed by allow `self` for iframes.

https://github.com/owncloud/ocis/pull/4031
https://github.com/owncloud/web/issues/7030
