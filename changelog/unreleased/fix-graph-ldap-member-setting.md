Bugfix: graph service did not honor the OCIS_LDAP_GROUP_SCHEMA_MEMBER setting

We fixed issue when using a custom LDAP attribute for group members. The graph service
did not honor the OCIS_LDAP_GROUP_SCHEMA_MEMBER environment variable

https://github.com/owncloud/ocis/issues/7032
