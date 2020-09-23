Bugfix: use NewNumericRangeInclusiveQuery for numeric literals

Some LDAP properties like `uidnumber` and `gidnumber` are numeric. When an OS tries to look up a user it will not only try to lookup the user by username, but also by the `uidnumber`: `(&(objectclass=posixAccount)(uidnumber=20000))`. The accounts backend for glauth was sending that as a string query `uid_number eq '20000'` and has been changed to send it as `uid_number eq 20000`. The removed quotes allow the parser in ocis-accounts to identify the numeric literal and use the NewNumericRangeInclusiveQuery instead of a TermQuery.

https://github.com/owncloud/ocis-glauth/issues/28
https://github.com/owncloud/ocis/accounts/pull/68
https://github.com/owncloud/ocis-glauth/pull/29
