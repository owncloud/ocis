Bugfix: Fix error code when a user can't disable a space

Previously, if the user couldn't disable a space due to wrong permissions, the
request returned a 404 error code, as if the space wasn't found even though
the space was visible. Now it will return the expected 403 error code.

https://github.com/owncloud/ocis/pull/11845
