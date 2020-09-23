Bugfix: exit when assets or config are not found

When a non-existing assets folder is specified, there was only a warning log statement and the service served
the builtin assets instead. It is safe to exit the service in such a scenario, instead of serving other assets
than specified. We changed the log level to `Fatal` on non-existing assets.
Similar for the web config, it was not failing on service level, but only showing an error in the web ui, wenn
the specified config file could not be found. We changed the log level to `Fatal` as well.

<https://github.com/owncloud/ocis/ocis-phoenix/pull/76>
<https://github.com/owncloud/ocis/ocis-phoenix/pull/77>
