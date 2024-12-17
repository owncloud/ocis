Bugfix: Fix possible race condition when a thumbnails is stored in the FS

A race condition could cause the thumbnail service to return a thumbnail
with 0 bytes or with partial content. In order to fix this, the service will
create a temporary file with the contents and then rename that file to its
final location.

https://github.com/owncloud/ocis/pull/10693
