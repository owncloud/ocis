Bugfix: Fix Password Reset

The `ocis idm resetpassword` always used the hardcoded `admin` name for the user. Now user name can be specified via the `--user-name` (`-u`) flag.

https://github.com/owncloud/ocis/pull/9479
