Bugfix: Fix default configuration for accessing shares

The storage provider mounted at `/home` should always have EnableHome set to `true`. The other storage providers should have it set to `false`.

https://github.com/owncloud/product/issues/205
https://github.com/owncloud/ocis-reva/pull/461


