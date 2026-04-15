Enhancement: Allow resetting IDM service user passwords

The `ocis idm resetpassword` command now supports a `--user-type` flag
to select the account type: `user` (default, ou=users) or `service`
(ou=sysusers). This allows resetting passwords for service accounts
(libregraph, idp, reva) which live in `ou=sysusers`. Previously, the
DN was hardcoded to `ou=users`, making it impossible to reset service
user passwords via the CLI.

https://github.com/owncloud/ocis/pull/12118
https://github.com/owncloud/ocis/issues/12106
