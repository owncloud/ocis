Bugfix: SpaceEditorWithoutTrashbin roles now correctly allow file editing

Fixed a bug where the *WithoutTrashbin space editor roles were rendered as read-only
in the Web frontend. The OCS PermissionWrite bit was not set for these roles because
the RoleFromResourcePermissions round-trip required RestoreRecycleItem, which these
roles intentionally omit.

https://github.com/owncloud/ocis/pull/12346
