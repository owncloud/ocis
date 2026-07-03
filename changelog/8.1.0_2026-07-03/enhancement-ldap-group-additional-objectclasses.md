Enhancement: Allow multiple objectClasses on group creation

Added support for configuring additional LDAP objectClasses when creating groups.
The new `OCIS_LDAP_GROUP_ADDITIONAL_OBJECTCLASSES` / `GRAPH_LDAP_GROUP_ADDITIONAL_OBJECTCLASSES`
environment variable accepts a list of extra objectClasses that are set alongside the
primary `GRAPH_LDAP_GROUP_OBJECTCLASS` when a new group is created in LDAP.

https://github.com/owncloud/ocis/pull/12229
