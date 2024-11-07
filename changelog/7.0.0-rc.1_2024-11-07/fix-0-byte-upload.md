Bugfix: Fix 0-byte file uploads

We fixed an issue where 0-byte files upload did not return the Location header.

https://github.com/owncloud/ocis/pull/10500
https://github.com/owncloud/ocis/issues/10469
