Bugfix: GetUserByClaim fixed for Active Directory

The reva ldap backend for the users and groups service did not hex escape
binary uuids in LDAP filter correctly this could cause problems in Active
Directory setups for services using the GetUserByClaim CS3 request with claim
"userid".

https://github.com/owncloud/ocis/pull/7476
https://github.com/owncloud/ocis/issues/7469
