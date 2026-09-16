Bugfix: Configure mount_id, events, and tokens for posix storage driver

Configure `mount_id`, `events`, and `tokens` in the POSIX driver configuration
for the data provider, and introduce `PosixNoEvents` for the storage provider.
Without `mount_id`, asynchronous file uploads on POSIX storage failed to
finalize because the `PostprocessingFinished` event check dropped the event as
belonging to a different storage provider, leaving uploads permanently stuck.

https://github.com/owncloud/ocis/pull/12780
