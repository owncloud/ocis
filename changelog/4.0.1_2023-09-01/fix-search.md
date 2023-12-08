Bugfix: Disallow sharee to search sharer files outside the share

When a file was shared with user(sharee) and the sharee searched the shared file the response contained unshared resources as well.

https://github.com/owncloud/ocis/pull/7184
