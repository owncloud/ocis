Enhancement: Cli to purge expired trash-bin items

Introduction of a new cli command to purge old trash-bin items.
The command is part of the `storage-users` service and can be used as follows:

`ocis storage-users trash-bin purge-expired`.

The `purge-expired` command configuration is done in the `ocis`configuration or as usual by using environment variables.

ENV `STORAGE_USERS_PURGE_TRASH_BIN_USER_ID` is used to obtain space trash-bin information and takes the system admin user as the default `OCIS_ADMIN_USER_ID`.
It should be noted, that this is only set by default in the single binary. The command only considers spaces to which the user has access and delete permission.

ENV `STORAGE_USERS_PURGE_TRASH_BIN_PERSONAL_DELETE_BEFORE` has a default value of `30 days`, which means the command will delete all files older than `30 days`.
The value is human-readable, valid values are `24h`, `60m`, `60s` etc. `0` is equivalent to disable and prevents the deletion of `personal space` trash-bin files.

ENV `STORAGE_USERS_PURGE_TRASH_BIN_PROJECT_DELETE_BEFORE` has a default value of `30 days`, which means the command will delete all files older than `30 days`.
The value is human-readable, valid values are `24h`, `60m`, `60s` etc. `0` is equivalent to disable and prevents the deletion of `project space` trash-bin files.

Likewise, only spaces of the type `project` and `personal` are taken into account.
Spaces of type `virtual`, for example, are ignored.

https://github.com/owncloud/ocis/pull/5500
https://github.com/owncloud/ocis/issues/5499
