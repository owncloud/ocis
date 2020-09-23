Bugfix: query numeric attribute values without quotes

Some LDAP properties like `uidnumber` and `gidnumber` are numeric. When an OS tries to look up a user it will not only try to lookup the user by username, but also by the `uidnumber`: `(&(objectclass=posixAccount)(uidnumber=20000))`. The accounts backend for glauth was sending that as a string query `uid_number eq '20000'` in the ListAccounts query. This PR changes that to `uid_number eq 20000`. The removed quotes allow the parser in ocis-accounts to identify the numeric literal.

<https://github.com/owncloud/ocis/glauth/issues/28>
<https://github.com/owncloud/ocis/glauth/pull/29>
<https://github.com/owncloud/ocis/accounts/pull/68>
