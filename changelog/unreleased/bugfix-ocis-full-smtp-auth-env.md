Bugfix: Pass SMTP password and authentication to ocis in the ocis_full example

The `ocis_full` deployment example defined `SMTP_PASSWORD` and
`SMTP_AUTHENTICATION` in its `.env` file but did not forward them to the `ocis`
container, so an SMTP server requiring authentication could not be configured.
`NOTIFICATIONS_SMTP_PASSWORD` and `NOTIFICATIONS_SMTP_AUTHENTICATION` are now
passed through.

https://github.com/owncloud/ocis/issues/11757
