Bugfix: Use O(1) document lookup instead of full search during reindexing

The `IndexSpace` bulk reindexer was using a full KQL search query per file
to check whether re-extraction was needed. On large indexes this query
took 600–950ms each, making a 61,000-file space take ~13.5 hours just to
walk. Replaced the per-file `Search()` call with an O(1) `Lookup()` using
Bleve's `DocIDQuery`, then comparing mtime and extraction status in memory.
This reduces per-file check time from ~800ms to <1ms.

https://github.com/owncloud/ocis/pull/12096
https://github.com/owncloud/ocis/issues/12093
