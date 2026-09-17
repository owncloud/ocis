Bugfix: Recover uploads stuck in processing when the blob commit fails

When the final blob commit failed after postprocessing, the node kept its
`processing` marker with no retry and no timeout. Every download then returned
`425 Too Early`, the file could not be deleted, and its reserved size kept
counting against the quota.

The commit is now retried with backoff. If it still fails, the node is reverted
to a recoverable failed state, which clears the processing marker and releases
the reserved quota.

https://github.com/owncloud/reva/pull/733
https://github.com/owncloud/ocis/pull/12970
