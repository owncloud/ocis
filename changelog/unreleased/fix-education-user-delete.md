Bugfix: Education user DELETE no longer 404s while leaving LDAP entry behind

The education user delete handler used `user.GetExternalID()` for the backend
DELETE, while the regular `/users` handler and the pre-v8.0 code path used
`user.GetId()`. With the default `RequireExternalID=false`, the LDAP backend
looked up the user by name-or-UUID, so the externalID never matched, the LDAP
entry was never removed, and the response was a 404. This is now fixed.

https://github.com/owncloud/ocis/pull/12405
