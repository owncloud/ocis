Bugfix: Serve index.html for directories

The static middleware in ocis-pkg now serves index.html instead of returning 404 on paths with a trailing `/`.

https://github.com/owncloud/ocis/pull/912
https://github.com/owncloud/ocis-pkg/issues/63
