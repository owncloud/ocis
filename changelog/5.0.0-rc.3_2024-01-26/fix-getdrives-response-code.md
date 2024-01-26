Bugfix: fix PATCH/DELETE status code for drives that don't support them

Updating and Deleting the virtual drives for shares is currently not supported. Instead
of returning a generic 500 status we return a 405 response now.

https://github.com/owncloud/ocis/pull/8235
https://github.com/owncloud/ocis/issues/7881
