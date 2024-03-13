Bugfix: Fix last month search

We've fixed the last month search edge case when currently is 31-th.

https://github.com/owncloud/ocis/issues/7629
https://github.com/owncloud/ocis/pull/7742

The issue is related to the build-in package behavior  https://github.com/golang/go/issues/31145
