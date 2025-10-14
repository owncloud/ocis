Enhancement: add named locks and refactor cache

Tags: ocis-pkg, accounts

We had the case that we needed kind of a named locking mechanism which enables us to lock only under certain conditions.
It's used in the indexer package where we do not need to lock everything, instead just lock the requested parts and differentiate between reads and writes.

This made it possible to entirely remove locks from the accounts service and move them to the ocis-pkg indexer.
Another part of this refactor was to make the cache atomic and write tests for it.

- remove locking from accounts service
- add sync package with named mutex
- add named locking to indexer
- move cache to sync package

https://github.com/owncloud/ocis/pull/1212
https://github.com/owncloud/ocis/issues/966
