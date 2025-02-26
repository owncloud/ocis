Bugfix: Fix status code for thumbnail requests

We fixed the status code returned by the thumbnails service when the image
source for a thumbnail exceeds the configured maximum dimensions or file size.
The service now returns a 403 Forbidden status code instead of a 500 Internal
Server Error status code.

https://github.com/owncloud/ocis/pull/10592
https://github.com/owncloud/ocis/issues/10589
