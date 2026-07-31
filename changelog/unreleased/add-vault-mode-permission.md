Enhancement: Add VaultMode permission

Add a new `VaultMode.ReadWriteEnabled` permission that gates the visibility of the
vault mode switcher in the web UI. The permission is assigned to the admin,
space admin and user roles. The user light role does not receive it.

https://github.com/owncloud/ocis/pull/12328
