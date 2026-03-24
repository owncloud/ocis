Bugfix: Always descend into directories during space reindexing

The search indexer's `IndexSpace` walk previously used `filepath.SkipDir`
to skip entire directory subtrees when the directory itself was already
indexed. After a failed or interrupted indexing run (e.g. Tika crash),
this caused thousands of unindexed files to be permanently skipped
because the parent directory's mtime had not changed. The indexer now
always descends into directories, relying on the O(1) per-file DocID
lookup to skip already-indexed files efficiently.

https://github.com/owncloud/ocis/pull/12119
