Bugfix: Remove group memberships when deleting a user

The LDAP backend in the graph API now takes care of removing a user's group
membership when deleting the user.

https://github.com/owncloud/ocis/issues/3027
