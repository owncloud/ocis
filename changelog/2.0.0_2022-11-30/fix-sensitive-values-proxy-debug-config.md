Bugfix: Fix make sensitive config values in the proxy's debug server

We've fixed a security issue of the proxy's debug server config report endpoint.
Previously sensitive configuration values haven't been masked. We now mask these values.

https://github.com/owncloud/ocis/pull/4086
