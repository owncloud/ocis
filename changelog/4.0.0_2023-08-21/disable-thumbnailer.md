Enhancement: Thumbnails can be disabled for webdav & web now

We added an env var `OCIS_DISABLE_PREVIEWS` to disable the thumbnails for web & webdav via a global setting.
For each service this behaviour can be disabled using the local env vars `WEB_OPTION_DISABLE_PREVIEWS` (old)
and `WEBDAV_DISABLE_PREVIEWS` (new).

https://github.com/owncloud/ocis/pull/6577
https://github.com/owncloud/ocis/issues/192