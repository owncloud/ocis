Bugfix: Thumbnail request limit

The `THUMBNAILS_MAX_CONCURRENT_REQUESTS` setting was not working correctly.
Previously it was just limiting the number of concurrent thumbnail downloads.
Now the limit is applied to the number thumbnail generations requests.
Additionally the webdav service is now returning a "Retry-After" header when
it is hitting the ratelimit of the thumbnail service.

https://github.com/owncloud/ocis/pull/10280
https://github.com/owncloud/ocis/pull/10225
