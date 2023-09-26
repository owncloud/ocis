Enhancement: New value `auto` for NOTIFICATIONS_SMTP_AUTHENTICATION

This cause the notifications service to automatically pick a suitable authentication
method to use with the configured SMTP server. This is also the new default behavior.
The previous default was to not use authentication at all.

https://github.com/owncloud/ocis/issues/7356
