Change: Provide cache for roles

In order to work efficiently with permissions we provide a cache for roles and a
middleware to update the cache based on roleIDs from the metadata context. It can be
used to check permissions in service handlers.

<https://github.com/owncloud/ocis-pkg/pull/59>
