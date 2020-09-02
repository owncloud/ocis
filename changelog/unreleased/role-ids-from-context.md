Change: Reuse roleIDs from the metadata context

The roleIDs of the authenticated user are coming in from the metadata context. Since we decided to move the role assignments over to the accounts service we need to start trusting those roleIDs from the metadata context instead of reloading them from disk on each request.

https://github.com/owncloud/ocis-settings/pull/69

