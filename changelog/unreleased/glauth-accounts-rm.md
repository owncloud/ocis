Change: The `glauth` and `accounts` services are removed

After switching the default configuration to libregraph/idm we could remove
the glauth and accounts services from the source code (they were already disabled
by default with the previous release)

https://github.com/owncloud/ocis/pull/3685
