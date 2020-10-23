Change: Rebuild index command for accounts

Tags: accounts

The index for the accounts service can now be rebuilt by running the cli command `./bin/ocis accounts rebuild`. It deletes all configured indices and rebuilds them from the documents found on storage.

https://github.com/owncloud/ocis/pull/748
