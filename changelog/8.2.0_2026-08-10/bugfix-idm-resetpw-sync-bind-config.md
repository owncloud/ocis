Bugfix: Keep bind config in sync when resetting a service user password

Resetting an IDM service user password with `ocis idm resetpassword
--user-type service` only changed the entry in the IDM directory. Services such
as auth-basic, users and groups still bound to LDAP as the `reva` service user
using the old `bind_password` from their configuration, so after a restart
those binds failed and regular admin login started returning 401 Unauthorized.

The command now rewrites the matching `bind_password` and
`service_user_passwords` keys in `ocis.yaml` when it is present, and always
prints the environment variables that must carry the new password for env-var
and distributed deployments it cannot rewrite.

https://github.com/owncloud/ocis/pull/12716
https://github.com/owncloud/ocis/issues/12569
