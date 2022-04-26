Bugfix: Thumbnails for `/dav/xxx?preview=1` requests

We've added the thumbnail rendering for `/dav/xxx?preview=1` requests, which was previously not supported because of missing routes. It now returns the same thumbnails as for
`/remote.php/dav/xxx?preview=1`.

We've also ensured that `/remote.php/webdav/xxx?preview=1` and `/webdav/xxx?preview=1` will be
routed to the correct service and always return a 404 Not Found, because Thumbnails are currently
not implemented for that route.

https://github.com/owncloud/ocis/pull/3567
