Change: Filesystem based index

Tags: accounts, storage

We replaced `bleve` with a new filesystem based index implementation. There is an `indexer` which is capable of
orchestrating different index types to build indices on documents by field. You can choose from the index types `unique`,
`non-unique` or `autoincrement`. Indices can be utilized to run search queries (full matches or globbing) on document
fields. The accounts service is using this index internally to run the search queries coming in via `ListAccounts` and
`ListGroups` and to generate UIDs for new accounts as well as GIDs for new groups.

The accounts service can be configured to store the index on the local FS / a NFS (`disk` implementation of the index)
or to use an arbitrary storage ( `cs3` implementation of the index). `cs3` is the new default, which is configured to
use the `metadata` storage.

https://github.com/owncloud/ocis/pull/709
