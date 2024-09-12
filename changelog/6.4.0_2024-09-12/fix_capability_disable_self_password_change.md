Bugfix: Set capability response `disable_self_password_change` correctly

The capability value `disable_self_password_change` was not being set correctly
when `user.passwordProfile` is configured as a read-only attribute.

https://github.com/owncloud/ocis/pull/9853
https://github.com/owncloud/enterprise/issues/6849
