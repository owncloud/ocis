Bugfix: Do not reindex a space twice at the same time

We fixed a problem where the search service reindexed a space while another
reindex process was still in progress.

https://github.com/owncloud/ocis/pull/5001
