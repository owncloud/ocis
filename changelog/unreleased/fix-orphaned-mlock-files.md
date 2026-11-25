Bugfix: Remove orphaned metadata files on delete and purge

The storage-users service left files on the disk when a user deleted a file or
purged it from the trash bin. Each orphaned file used an inode. On a volume
with a hard inode limit, the volume can become full.

The delete operation now removes the `.mlock` file of the node. A purge now
removes the `.mlock` files of the node and of its revisions. A purge from the
trash bin now finds the revisions of the node. Before, the revision files,
their metadata files and their blobs stayed on the disk.

This fix does not remove files that are already orphaned.

https://github.com/owncloud/ocis/pull/13563
https://github.com/owncloud/ocis/issues/13476
https://github.com/owncloud/reva/pull/849
