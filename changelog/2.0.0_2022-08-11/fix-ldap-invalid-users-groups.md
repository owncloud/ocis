Bugfix: Fix handling of invalid LDAP users and groups

We fixed an issue where ocis would exit with a panic when LDAP users
or groups where missing required attributes (e.g. the id)

https://github.com/owncloud/ocis/issues/4274
