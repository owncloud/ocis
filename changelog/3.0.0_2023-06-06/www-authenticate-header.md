Bugfix: Fix authenticate headers for API requests

We changed the www-authenticate header which should not be sent when the `XMLHttpRequest` header is set.

https://github.com/owncloud/ocis/pull/5992
https://github.com/owncloud/ocis/issues/5986
