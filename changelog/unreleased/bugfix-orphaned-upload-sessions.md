Bugfix: Release the quota of upload sessions with unreadable node metadata

When an upload's target node lost its metadata, e.g. because an ancestor was
moved to the trash while the upload was still in flight, the upload could never
finish processing. It stayed in "Processing" forever, could not be downloaded or
deleted, and kept consuming the space quota.

Cleaning such a session up with `ocis storage-users uploads sessions --clean`
did not help either: it removed the uploaded bytes and the session info file
before failing to revert the node, so it destroyed the only copy of the data
without releasing any quota.

Cleanup now reverts the node before removing anything irreversible and falls
back to the session metadata when the node cannot be read, so the quota is
released and the orphaned node is removed. If the quota cannot be released the
upload is kept so it can be retried instead of being lost.

A new `--orphaned` filter lists the affected sessions:

```
ocis storage-users uploads sessions --orphaned
ocis storage-users uploads sessions --orphaned --clean
```

Note that evaluating the filter reads the node metadata of every session, so it
is only done when the flag is set.

https://github.com/owncloud/ocis/pull/12739
