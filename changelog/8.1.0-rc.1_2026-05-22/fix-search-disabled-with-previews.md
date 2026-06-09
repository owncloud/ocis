Bugfix: Search no longer disabled when OCIS_DISABLE_PREVIEWS=true

Setting OCIS_DISABLE_PREVIEWS=true removed the WebDAV REPORT routes from
the router, breaking search on /dav/files, /dav/spaces and /webdav. The
search routes are now registered independently of the preview flag.

https://github.com/owncloud/ocis/pull/12303
