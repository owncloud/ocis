Bugfix: ignore case when comparing objectclass values

The LDAP equality comparison is specified as case insensitive. We fixed the comparison for objectclass properties.

https://github.com/owncloud/ocis-glauth/pull/26
