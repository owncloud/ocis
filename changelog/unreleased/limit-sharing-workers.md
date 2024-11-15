Bugfix: We now limit the number of workers of the jsoncs3 share manager

We now restrict the number of workers that look up shares to 5. The number can be changed with `SHARING_USER_JSONCS3_MAX_CONCURRENCY`.

https://github.com/owncloud/ocis/pull/10552
