Bugfix: Fix error logging when there is no thumbnail for a file

We've fixed the behavior of the logging when there is no thumbnail for a file
(because the filetype is not supported for thumbnail generation).
Previously the WebDAV service always issues an error log in this case. Now, we don't log this event any more.

https://github.com/owncloud/ocis/pull/2702
