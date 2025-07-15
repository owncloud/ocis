Bugfix: Trigger a rescan of spaces in the search index when items have changed

The search service now scans spaces when items have been changed. This fixes the problem
that mtime and treesize propagation was not reflected in the search index properly.

https://github.com/owncloud/ocis/pull/4777
https://github.com/owncloud/ocis/issues/4410
