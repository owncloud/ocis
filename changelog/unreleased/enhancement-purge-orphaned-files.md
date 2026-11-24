Enhancement: Add a command to purge orphaned files

We added the command `ocis trash purge-orphaned-files`. Older versions of oCIS
left files on the disk when a user deleted a file or purged it from the trash
bin. Each orphaned file used an inode.

The command removes a lock file when its node does not exist. It also removes
a revision with its metadata file, its lock file and its blob, when its node
does not exist and is not in the trash bin. By default, the command only
prints the files. It does not remove files that are younger than one hour.

https://github.com/owncloud/ocis/pull/13570
https://github.com/owncloud/ocis/issues/13476
