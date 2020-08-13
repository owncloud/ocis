Bugfix: exit when assets are not found

When a non-existing assets folders is specified, there was only a warning log statement and the service served
the builtin assets instead. It is safe to exit the service in such a scenario, instead of serving other assets
than specified. We changed the log level to `Fatal` on non-existing assets.

https://github.com/owncloud/ocis-phoenix/pull/76
