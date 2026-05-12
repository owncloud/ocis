Enhancement: Add vault storage with MFA-protected access

Added a dedicated vault storage that can be protected with MFA. A separate
`storage-users` service instance configured in vault mode runs and serves
`/vault/users` and `/vault/projects` mount points with a dedicated
`VaultStorageProviderID`. The `graph` service gained a new vault mode
(`GRAPH_ENABLE_VAULT_MODE`) that serves the vault API under the `/vault`
prefix. The storage registry now routes vault-specific requests exclusively to
the vault storage provider, preventing accidental access to vault spaces when
no explicit storage ID is provided.

MFA status is propagated through gRPC metadata
and forwarded in HTTP headers for WOPI/collaboration flows.

https://github.com/owncloud/ocis/pull/12108
