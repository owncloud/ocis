Bugfix: Fix share metadata corruption during concurrent share operations

When multiple sharing service replicas processed share operations concurrently
for the same user, the share metadata could become corrupted with references to
missing data, making all shares inaccessible to that user. The received share
cache now uses compare-and-swap (etag) validation to detect concurrent writes
and retries gracefully, preventing metadata corruption.

https://github.com/owncloud/ocis/pull/12621
