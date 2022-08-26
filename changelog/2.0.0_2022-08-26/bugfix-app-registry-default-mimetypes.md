Bugfix: Fix configuration of mimetypes for the app registry

We've fixed the configuration option for mimetypes in the app registry.
Previously the default config would always be merged over the user provided
configuration. Now the default mimetype configuration is only used if the user does not
providy any mimetype configuration (like it is already done in the proxy with the routes configuration).

https://github.com/owncloud/ocis/pull/4411
