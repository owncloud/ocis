Change: Rebuild index command for accounts

Tags: accounts

The index for the accounts service can now be rebuilt by running the cli command `./bin/ocis accounts rebuild`.
It deletes all configured indices and rebuilds them from the documents found on storage. For this we also introduced
a `LoadAccounts` and `LoadGroups` function on storage for loading all existing documents.

https://github.com/owncloud/ocis/pull/748
