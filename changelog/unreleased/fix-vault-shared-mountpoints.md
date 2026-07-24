Bugfix: Show vault-shared mountpoints in the vault drive list

Shares of vault resources are mountpoints hosted on the shares storage provider
that grant into the vault storage provider. When listing spaces with the vault
`storage_id`, the storage registry skipped the shares provider entirely, so the
vault share mountpoint was missing from the vault drive list while still showing
up in the regular drive list.

The registry now also queries the shares provider when the vault storage id is
requested and segregates share mountpoints by their `grantStorageID`, so a vault
share only appears in the vault drive list and a regular share only in the
regular one.

https://github.com/owncloud/ocis/pull/12644
https://github.com/owncloud/reva/pull/666
