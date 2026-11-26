Bugfix: Cache decoded node attributes

The file metadata cache keeps a bounded, process-local layer of already decoded
node attributes, so reading cached attributes no longer decodes them on every
access. Every read still checks the value in the configured cache store, so
the result is the same as before. By default the layer is limited to 16384
entries and about 64 MiB per cache.

https://github.com/owncloud/ocis/pull/13597
https://github.com/owncloud/reva/pull/864
