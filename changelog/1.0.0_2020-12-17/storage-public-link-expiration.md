Bugfix: Disable public link expiration by default

Tags: storage

The public link expiration was enabled by default and didn't have a default expiration span by default, which resulted in already expired public links coming from the public link quick action. We fixed this by disabling the public link expiration by default.

https://github.com/owncloud/ocis/issues/987
https://github.com/owncloud/ocis/pull/1035
