Bugfix: Inconsistency env var naming for LDAP filter configuration

There was a naming inconsitency for the enviroment variables used to define
LDAP filters for user and groups queries. Some services used `LDAP_USER_FILTER`
while others used `LDAP_USERFILTER`. This is now changed to use `LDAP_USER_FILTER`
and `LDAP_GROUP_FILTER`.

Note: If your oCIS setup is using an LDAP configuration that has any of the
`*_LDAP_USERFILTER` or `*_LDAP_GROUPFILTER` environment variables set, please
update the configuration to use the new unified names `*_LDAP_USER_FILTER`
respectively `*_LDAP_GROUP_FILTER` instead.

https://github.com/owncloud/ocis/issues/3890
