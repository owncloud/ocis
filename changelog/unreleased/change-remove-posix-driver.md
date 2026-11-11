Change: Remove the posix driver

We removed the deprecated `posix` driver from oCIS.

The service `storage-users` no longer accepts the value `posix` for its driver.
The related configuration section and environment variables are removed. The
documentation for the posix driver is also removed.

The commands `ocis decomposedfs`, `ocis migrate`, and `ocis revisions` now use
the decomposedfs `timemanager` and `lookup` packages.

https://github.com/owncloud/ocis/pull/13187
