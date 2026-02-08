Bugfix: Make tag unassignment idempotent

The DELETE tags endpoint now returns success when the requested tag is already absent from the file's metadata, instead of returning HTTP 400 with a misleading error message. The TagsRemoved event is always published so the search index stays in sync even when file metadata and the search index are out of sync.

https://github.com/owncloud/ocis/pull/12001
