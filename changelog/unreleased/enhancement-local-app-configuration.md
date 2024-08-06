Enhancement: Local WEB App configuration

We've added a new feature which allows configuring applications individually instead of using the global apps.yaml file.
With that, each application can have its own configuration file, which will be loaded by the WEB service.

The local configuration has the highest priority and will override the global configuration.
The Following order of precedence is used: local.config > global.config > manifest.config.

Besides the configuration, the application now be disabled by setting the `disabled` field to `true` in one of the configuration files.

https://github.com/owncloud/ocis/pull/9691
https://github.com/owncloud/ocis/issues/9687
