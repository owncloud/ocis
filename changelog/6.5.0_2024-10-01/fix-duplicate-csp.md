Bugfix: Remove duplicate CSP header from responses

The web service was adding a CSP on its own, and that one has been removed. The proxy service will take care of the CSP header.

https://github.com/owncloud/ocis/pull/10146
