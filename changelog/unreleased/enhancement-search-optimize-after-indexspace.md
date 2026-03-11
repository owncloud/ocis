Enhancement: Optimize search index after bulk reindexing

After an `IndexSpace` walk completes, the search engine now triggers a
segment merge (compaction) on the bleve index. Over time, writes create
multiple index segments that degrade query performance. The new
`Optimize()` method calls bleve's `ForceMerge` to consolidate all
segments into one, improving subsequent search and lookup speed. This is
especially beneficial after bulk reindexing large spaces.

https://github.com/owncloud/ocis/pull/12104
https://github.com/owncloud/ocis/issues/12093
