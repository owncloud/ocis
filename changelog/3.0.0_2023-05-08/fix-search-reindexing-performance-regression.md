Bugfix: Fix Search reindexing performance regression

We've fixed a regression in the search service reindexing step, causing the
whole space to be reindexed instead of just the changed resources.

https://github.com/owncloud/ocis/pull/6085
