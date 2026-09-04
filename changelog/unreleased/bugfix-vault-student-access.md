Bugfix: Enforce the vault permission on the Safe/Vault routes

Access to the Safe/Vault is now authorized on the per-user vault permission. The
`/vault/graph` routes return 403 for users without the permission, and the web
router routes them to the access-denied page.

https://github.com/owncloud/ocis/pull/12881
