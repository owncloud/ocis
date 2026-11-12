Enhancement: Add a command to remove the trash bin access from existing space viewers

The "Can view" space role no longer allows listing the trash bin of the space.
The permissions of a space member are stored when the member is added, so
members that got the "Can view" role before this change could still list the
trash bin, and the members list showed them with "Custom permissions".

The new command `ocis storage-users spaces fix-viewer-trash` removes the trash
bin access from these members. It runs as a dry run by default and only changes
members whose permissions are exactly the ones of the old "Can view" role.

https://github.com/owncloud/ocis/pull/13393
