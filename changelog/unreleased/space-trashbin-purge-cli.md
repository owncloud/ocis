Enhancement: Cli to purge expired trash-bin items

Introduction of a new cli command to purge old trash-bin items.
The command is part of the `storage-users` service and can be used as follows:

`storage-users trash-bin purge-expired --purge-before=24h --user-id=some-user-id`.

The `purge-before flag` has a default value of `30 days`, which means the command will delete all files older than `30 days`.
The value is human-readable, valid values are `24h`, `60m`, `60s` etc.

The `user-id flag` takes the system admin user as the default. It should be noted, that this is only set by default in the single binary.
The command only considers spaces to which the user has access and delete permission.

Likewise, only spaces of the type `project` and `personal` are taken into account.
Spaces of type virtual`, for example, are ignored.

https://github.com/owncloud/ocis/pull/5500
https://github.com/owncloud/ocis/issues/5499
