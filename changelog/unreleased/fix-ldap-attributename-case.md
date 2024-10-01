Bugfix: always treat LDAP attribute names case-insensitively

We fixes a bug where some LDAP attributes (e.g. owncloudUUID) were not
treated case-insensitively.

https://github.com/owncloud/ocis/pull/10204
https://github.com/owncloud/ocis/issues/10200
