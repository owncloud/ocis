Bugfix: Make the share cache locks cancellable

Requests that wait for a lock of the share caches of the jsoncs3 share manager
now give up when their context is done, instead of waiting until they get the
lock. Requests that the client already abandoned no longer pile up and keep
their memory under load. The provider cache also rejects updates with an empty
storage, space or share id.

https://github.com/owncloud/ocis/pull/13595
https://github.com/owncloud/reva/pull/838
https://github.com/owncloud/reva/pull/847
