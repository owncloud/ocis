Bugfix: Stop logging "no user in context" for data gateway requests

The proxy's create-home middleware logged an error for every request that
carried a reva token but no user. reva's internal downloads through the data
gateway (`/data`) are such requests, so every file download added an
error-level line. These requests are now passed on with a debug-level log.

https://github.com/owncloud/ocis/issues/13101
