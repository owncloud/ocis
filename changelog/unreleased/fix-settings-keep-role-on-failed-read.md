Bugfix: Keep the role of a user when the settings service cannot read it

The settings service could remove the role of a user without an error. When
the read of a role assignment failed, the service returned an empty role list.
The proxy then assigned the default role "User", and the settings service
deleted the old role, for example Admin.

We fixed these causes:

- The self-assignment of the default role failed open. When the settings
  service could not read the existing assignments, it wrote the default role.
  It now returns an error and does not write.
- Each settings pod has its own cache. After a role change on one pod, the
  cached listing of the account folder on another pod named an assignment
  that was gone. The settings service now drops that listing and reads the
  folder again.
- The list of the assignments of a role was empty when one account folder was
  gone. The settings service now skips only that account.

https://github.com/owncloud/ocis/pull/13526
https://github.com/owncloud/reva/pull/820
