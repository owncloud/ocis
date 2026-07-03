Enhancement: Add `ocis search optimize` CLI command

Added a new `ocis search optimize` command that compacts the search index
by merging Bleve segments, without re-indexing content. The command opens
the index directly (without requiring the search service to be running),
making it safe to run during maintenance windows without blocking search
queries.

This is useful after bulk reindexing operations that create many small
index segments, which can degrade search performance over time.

https://github.com/owncloud/ocis/pull/12136
