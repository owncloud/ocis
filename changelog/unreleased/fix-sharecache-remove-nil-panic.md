Bugfix: Fix a nil pointer panic when a share is removed

We updated reva to fix a nil pointer panic in the jsoncs3 share cache. The
panic occurred when a user had no cache entry for the space of the removed
share. The panic stopped the whole process. All removal paths were affected,
for example unsharing and the cleanup of orphaned shares.

https://github.com/owncloud/ocis/pull/13472
https://github.com/owncloud/reva/pull/814
