Enhancement: allow to maintain the last sign-in timestamp of a user

When the LDAP identity backend is configured to have write access to the database
we're now able to maintain the ocLastSignInTimestamp attribute for the users.

This attribute is return in the 'signinActivity/lastSuccessfulSignInDateTime'
properity of the user objects. It is also possible to $filter on this attribute.

Use e.g. '$filter=signinActivity/lastSuccessfulSignInDateTime le 2023-12-31T00:00:00Z'
to search for users that have not signed in since 2023-12-31.
Note: To use this type of filter the underlying LDAP server must support the
'<=' filter. Which is currently not the case of the built-in LDAP server (idm).

https://github.com/owncloud/ocis/pull/9942
https://github.com/owncloud/ocis/pull/10111
