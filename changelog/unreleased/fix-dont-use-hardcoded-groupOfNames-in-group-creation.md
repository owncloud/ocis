Bugfix: Don't use hardcoded groupOfNames in group creation

When creating a group with different objectClass, it will always use groupOfNames instead of the one provided in the config.
The server now creates groups using the objectClass defined in the config.

https://github.com/owncloud/ocis/pull/11776
