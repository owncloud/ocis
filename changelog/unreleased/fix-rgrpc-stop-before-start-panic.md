Bugfix: Fix a panic when a service stops during its start

We updated reva to fix a panic in the reva grpc servers. A service panicked
when it got a stop signal before its grpc server was ready. This occurred, for
example, when the sharing service waited for NATS during its start. The panic
replaced the normal shutdown with a crash. Now the service stops normally.

https://github.com/owncloud/ocis/pull/13472
https://github.com/owncloud/reva/pull/814
