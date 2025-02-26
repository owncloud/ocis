Bugfix: The user attributes `userType` and `memberOf` are readonly

The graph API now treats the user attributes `userType` and `memberOf` as
read-only. They are not meant be updated directly by the client.

https://github.com/owncloud/ocis/pull/9867
https://github.com/owncloud/ocis/issues/9858
