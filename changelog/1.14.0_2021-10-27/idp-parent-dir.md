Bugfix: Create parent directories for idp configuration

The parent directories of the identifier-registration.yaml config file might
not exist when starting idp. Create them, when that is the case.

https://github.com/owncloud/ocis/issues/2667
