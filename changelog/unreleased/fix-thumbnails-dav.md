Bugfix: Thumbnails for `/dav/xxx?preview=1` requests

We've added the thumbnail rendering for `/dav/xxx?preview=1`, `/remote.php/webdav/{relative path}?preview=1` and `/webdav/{relative path}?preview=1` requests, which was previously not supported because of missing routes. It now returns the same thumbnails as for
`/remote.php/dav/xxx?preview=1`.

https://github.com/owncloud/ocis/pull/3567
