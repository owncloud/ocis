Bugfix: Store user passwords hashed in idm

Support for hashing user passwords was added to libregraph/idm. The graph API will
now set userpasswords using the LDAP Modify Extended Operation (RFC3062). In the default
configuration passwords will be hashed using the argon2id algorithm.

https://github.com/owncloud/ocis/issues/3778
https://github.com/owncloud/ocis/pull/4053
