Bugfix: Fix inverted PROPFIND public link resolution flag

The ocdav PROPFIND handler resolves public link shares to populate the
oc:share-type property. A new toggle for skipping that lookup was wired up
backwards, so setting it enabled the extra lookup instead of skipping it. The
condition has been fixed and the toggle is now exposed as
OCDAV_DISABLE_PROPFIND_PUBLIC_LINK_RESOLUTION, which can be set to reduce load
on services for large collections.

https://github.com/owncloud/ocis/pull/13020
