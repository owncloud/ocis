Enhancement: Add a command to purge orphaned files

We added the command `ocis trash purge-orphaned-files`. Older versions of oCIS
left files on the disk when a user deleted a file or purged it from the trash
bin. Each orphaned file used an inode.

For a node in the trash bin, the command removes the lock files of its
revisions. For a node that does not exist and is not in the trash bin, it
removes the lock file of the node, and the revisions with their metadata
files, their lock files and their blobs. By default, the command only
prints a summary. The `--verbose` option prints each file. The command does
not remove files that are younger than one hour.

https://github.com/owncloud/ocis/pull/13570
https://github.com/owncloud/ocis/issues/13476
