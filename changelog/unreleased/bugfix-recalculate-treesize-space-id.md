Bugfix: Accept CS3 space ids and fail properly in recalculate-treesize

The `ocis storage-users spaces recalculate-treesize` command rejected the CS3
space id as exposed by the graph API, WebDAV and the web UI, and exited with 0
and printed `All treesizes are correct` even when a space could not be
processed. The `--space-id` flag now accepts `<storage-id>$<space-id>!<opaque-id>`,
`<storage-id>$<space-id>` and the bare `<space-id>`. If any space fails, the
remaining spaces are still processed, but the command exits non-zero with a
summary of the failed spaces. Unreadable nodes, no spaces found in the storage
root, and space ids that would escape the storage root are errors as well, and a
treesize is never written from an incomplete sum. The documentation now also states that the
`storage-users` service must be stopped while the command runs, because it
serves treesizes from the file metadata cache.

https://github.com/owncloud/ocis/issues/13275
https://github.com/owncloud/ocis/pull/13409
