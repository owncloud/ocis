Bugfix: Fixes "unaligned 64-bit atomic operation" panic on 32-bit ARM

sync/cache had uint64s that were not 64-bit aligned causing panics
on 32-bit systems during atomic access

https://github.com/owncloud/ocis/pull/1888
https://github.com/owncloud/ocis/issues/1887
