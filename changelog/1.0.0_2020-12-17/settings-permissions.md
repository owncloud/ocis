Bugfix: Permission checks for settings write access

Tags: settings

There were several endpoints with write access to the settings service that were not protected by permission checks. We introduced a generic settings management permission to fix this for now. Will be more fine grained later on.

https://github.com/owncloud/ocis/pull/1092
