Enhancement: Use non root user for the owncloud/ocis docker image

The owncloud/ocis docker image now uses a non root user and enables you to set a different user with the docker `--user` parameter. The default user has the UID 1000 is part of a group with the GID 1000.

This is a breaking change for existing docker deployments. The permission on the files and folders in persistent volumes need to be changed to the UID and GID used for oCIS (default 1000:1000 if not changed by the user).

https://github.com/owncloud/ocis/pull/2380
