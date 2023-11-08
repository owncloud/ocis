Bugfix: Deprecate redundant encryptions settings for notification service

The values `tls` and `ssl` for the `smtp_encryption` configuration setting are
duplicates of `starttls` and `ssltls`. They have been marked as deprecated.
A warning will be logged when they are still used. Please use `starttls` instead
for `tls` and `ssltls` instead of `ssl.

https://github.com/owncloud/ocis/issues/7345
