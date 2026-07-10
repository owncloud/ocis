Bugfix: Pass the space ID when purging revision blobs

The `revisions purge` command removed a revision's metadata but did not delete
its blob, because it called the blobstore without the space ID. The blobstore
builds the blob path from the space ID and blob ID, so an empty space ID
targeted the wrong path: the deletion was a no-op (S3) or missed the file
(POSIX), leaving the blob orphaned while the revision was already removed. The
space ID parsed from the revision path is now passed to the blobstore.

https://github.com/owncloud/ocis/issues/11644
https://github.com/owncloud/ocis/pull/12422
