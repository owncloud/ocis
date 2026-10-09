Bugfix: Prevent lost notifications from a concurrent userlog event list race

The userlog service processes events with `MaxConcurrency` workers. Updating a
user's stored event list (`alterUserEventList`) and the shared global events
(`alterGlobalEvents`) was a read-modify-write against the key-value store with
no locking, so two workers operating on the same user (or on the single
global-events key) could interleave their read and write and silently drop each
other's update — losing in-app notifications. The read-modify-write cycles are
now serialized with a mutex, mirroring the activitylog service.

https://github.com/owncloud/ocis/pull/12415
