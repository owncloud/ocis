Bugfix: Open Safe permanent links from outside the Safe

A permanent link to a file in the Safe, opened in a regular Drive tab, showed
"An error occurred while resolving the private link" instead of the file. The
Web frontend recognized Safe files by a storage provider id that the server
only sends to sessions already inside the Safe, so a Drive tab never redirected
the link into the Safe.

The Safe storage provider id is fixed, so the frontend now uses it as a
constant, like the share jail and OCM ids. This also keeps Safe notifications
out of the Drive view and lets the invite form recognize Safe resources in
Drive mode. The `vaultStorageProvider` getter of the web-pkg capability store
is removed; use `VAULT_STORAGE_PROVIDER_ID` from web-client instead.

https://github.com/owncloud/ocis/pull/13001
