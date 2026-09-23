Bugfix: Resolve a permanent link to a vault resource from outside the vault

A permanent link copied inside the vault always points at the plain
`<server>/f/<fileId>` form, without the `/vault` path prefix. Opening it
therefore boots the Web frontend in regular drive mode, which is supposed to
notice that the file lives in the vault and reload itself into the vault scope.
That check compares the storage provider of the file id against the
`vault_storage_provider` capability, and the server only filled that capability
in for requests carrying `?vault=true`. A drive mode session was told
`vault.enabled: true` with an empty provider id, the comparison never matched,
the reload never happened, and the page then tried to read a vault resource with
drive mode clients that cannot reach it.

The vault storage provider id is a fixed constant rather than a deployment
secret, so it is now announced on every capabilities response. Besides the
permanent link this also repairs two other checks that ran against the empty
value outside the vault: vault notifications are filtered out of the regular
view again, and the invite form recognizes vault resources.

https://github.com/owncloud/ocis/pull/12989
