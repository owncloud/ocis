Bugfix: Thumbnail request limit

The `THUMBNAILS_MAX_CONCURRENT_REQUESTS` setting was not working correctly.
Previously it was just limiting the number of concurrent thumbnail downloads.
Now the limit is applied to the number thumbnail generations requests.

https://github.com/owncloud/ocis/pull/10225
