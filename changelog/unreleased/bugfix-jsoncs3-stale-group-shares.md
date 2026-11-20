Bugfix: Remove stale group shares in the share manager

A stale group share is a share ID in the group cache that has no share. On
each listing of the received shares, the share manager found the stale ID
again and uploaded the received.json file of the user again with the same
content. The stale ID stayed, so this repeated on each listing.

The share manager now removes a stale group share from the group cache. It
does this only when it read the share file of the space. It reads each share
from the same copy of that file, so a parallel listing cannot make valid
shares look stale. A Remove no longer uploads a file when the share is not in
it. A listing of the created shares now also reads each share from one copy
of the share file. Before, a parallel listing could hide some shares of a
space from the creator.

An update of a share now reads the share again under the space lock. Before,
a parallel failed upload in the same space could remove the space from
memory. The update then reported success, but it did not write the change.
An update now also returns NotFound for a share that a parallel unshare
removed. Before, the update could write the removed share again.

When the upload of a Remove fails, the share manager now reads the file again
on the next call. Before, it kept the share out of memory, but the share stayed
in the file. A retry then found nothing to remove, did not upload, and reported
success. In the received.json file of a user, this also occurred inside one
Remove after a temporary upload error. The same applies to an Add. Before, a
failed upload of an Add kept the new entry in memory. For example, a client got
an error when it accepted a share, but this pod showed the share as accepted.

An unshare and an update of a share now read the share file of the space
again when the share is not in memory. Before, a listing that got a temporary
NotFound left an empty copy of the space in memory. An unshare then did not
remove the share from the file, but it reported success. An update returned
NotFound for a share that existed.

The share manager removes a stale entry only after a grace period of 5
minutes. During a rolling upgrade, an older version can write the grantee
entry before the share. Without the grace period, a listing in this window
removed the entry, and no group member saw the share. A listing that finds
the share again restarts the grace period. Thus two short losses do not add
up to one long loss. Each pod counts the grace period itself. A pod remembers
a stale entry for 1 hour after the last listing that found it. Thus the
grace period also ends on a pod that finds the entry only once in some
minutes. A pod looks for the entries to forget at most once a minute. Thus a
listing that finds many stale entries does not hold the lock for a long time.
When the removal of a stale entry fails, the next listing tries the removal
again at once. The pod forgets the entry only after the removal.

The share manager now also removes a stale entry when the share file of the
space does not exist. A failed first share in a space leaves such an entry.
A download that gets NotFound is not enough proof, because the NotFound can
be temporary. The share manager checks the file a second time.

The share manager now writes a share before the entry for the grantee. When a
write fails, it removes the entries that it wrote, also the entry for the
grantee. The removal also runs when the request was cancelled. Each removal
has its own timeout, so a slow removal does not stop the next one. Thus a
failed share does not block a new share with the same grantee. The share
manager removes the share and the entry of the creator at the same time, so a
failed share waits for at most two timeouts. An unshare
removes the share also when the entry for the grantee cannot be removed, and
it reports success. Before, the client got an error for a share that was
gone. The next listing of the grantee removes the stale entry.

The share manager now removes many shares together, with one upload for each
document. A purge of a space with 20 group shares did 41 uploads and now does
3. A listing that finds 20 expired shares did 60 uploads and now does 3. A
listing removes the expired shares of all spaces together. The setting
SHARING_USER_JSONCS3_MAX_CONCURRENCY limits the parallel removals. Its default
does not change and is 1. The removals always use a limit of at least 3. Before,
the removal of one share updated its 3 documents in parallel. With a limit of
1, an unshare needed one more round trip to the storage.

The share cleanup of the ocis shares commands now removes all orphaned shares
together, with one upload for each document. Before, it removed the shares one
at a time. 20 orphaned shares in one space did 60 uploads and now do 3. A Stat
error no longer stops the cleanup with a panic. The cleanup skips that share
and checks the next share. The cleanup also stops with an error log when it
cannot get a gateway client. Before, it stopped with a panic.

The share manager now sends one ShareExpired event for each expired share.
Before, each listing that found the expired share sent an event, also when
the removal failed. Thus the users got many expiry notifications.

The share manager no longer keeps an empty space entry in the created cache
and in the group cache. Before, each listing read the share file of a
deleted space again.

When a space is deleted, the share manager now logs the errors of the share
cleanup. Before, it sent them to a disabled logger.

Each log line of a share listing now has the fields of one share only. Before,
a later line had the fields of all earlier shares, and a log search found the
wrong share.

A Remove on the created cache or on the group cache no longer panics when the
user or group has no entry for the space of the share, or when the share ID
has no delimiter. Before, the panic stopped the sharing service, also in an
unshare and in the cleanup of orphaned shares.

https://github.com/owncloud/ocis/pull/13544
https://github.com/owncloud/reva/pull/832
https://github.com/owncloud/reva/pull/813
