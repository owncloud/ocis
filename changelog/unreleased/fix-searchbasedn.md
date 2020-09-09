Bugfix: Use searchBaseDN if already a user/group name

In case of the searchBaseDN already referencing a user or group, the search query was ignoring the user/group name entirely, because the searchBaseDN is not part of the LDAP filters. We fixed this by including an additional query part if the searchBaseDN contains a CN.

https://github.com/owncloud/product/issues/214
https://github.com/owncloud/ocis-glauth/pull/32
