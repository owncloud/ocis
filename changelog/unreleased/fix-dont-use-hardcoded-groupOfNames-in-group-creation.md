Bugfix: Don't use hardcoded groupOfNames in group creation

Formerly, when creating a group with a different objectClass, it will always use groupOfNames instead of the one provided in the config.
Now, the server creates groups using the objectClass defined in the config.

https://github.com/owncloud/ocis/pull/11776
