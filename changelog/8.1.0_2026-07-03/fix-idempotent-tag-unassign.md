Bugfix: Make tag unassignment idempotent and handle publish failures

The DELETE tags endpoint now returns success when the requested tag is already absent from the file's metadata, instead of returning HTTP 400 with a misleading error message. The TagsRemoved event is always published so the search index stays in sync even when file metadata and the search index are out of sync. If event publishing fails, the metadata change is rolled back and HTTP 500 is returned to avoid leaving the system in an inconsistent state.

https://github.com/owncloud/ocis/pull/12001
