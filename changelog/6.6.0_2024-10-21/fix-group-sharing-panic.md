Bugfix: Fix panic when sharing with groups

We fixed a bug which caused a panic when sharing with groups, this only happened under a heavy load.
Besides the bugfix, we also reduced the number of share auto accept log messages to avoid flooding the logs.

https://github.com/owncloud/ocis/pull/10279
https://github.com/owncloud/ocis/issues/10258
