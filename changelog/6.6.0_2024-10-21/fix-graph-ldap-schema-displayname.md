Bugfix: graph service now supports `OCIS_LDAP_USER_SCHEMA_DISPLAYNAME` env var

To align with the other services the graph service now supports the
`OCIS_LDAP_USER_SCHEMA_DISPLAYNAME` environment variable to configure the LDAP
attribute that is used for display name attribute of users.

`LDAP_USER_SCHEMA_DISPLAY_NAME` is now deprecated and will be removed in a future
release.

https://github.com/owncloud/ocis/issues/10257
