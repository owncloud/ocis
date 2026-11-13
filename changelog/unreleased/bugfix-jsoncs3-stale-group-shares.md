Bugfix: Remove stale group shares in the share manager

A stale group share is a share ID in the group cache that has no share. On
each listing of the received shares, the share manager found the stale ID
again and uploaded the received.json file of the user again with the same
content. The stale ID stayed, so this repeated on each listing.

The share manager now removes a stale group share from the group cache. It
does this only when it read the share file of the space. A Remove no longer
uploads a file when the share is not in it.

The share manager now writes a share before the entry for the grantee, and it
removes the entries that it wrote when a write fails. An unshare removes the
share also when the entry for the grantee cannot be removed. The next listing
of the grantee removes the stale entry.

The share manager now removes many shares together, with one upload for each
document. A purge of a space with 20 group shares did 41 uploads and now does
3. A listing that finds 20 expired shares did 60 uploads and now does 3.

When a space is deleted, the share manager now logs the errors of the share
cleanup. Before, it sent them to a disabled logger.

https://github.com/owncloud/ocis/pull/13460
https://github.com/owncloud/reva/pull/813
