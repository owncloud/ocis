Enhancement: Unify LDAP config settings accross services

The storage services where updated to adapt for the recent changes of the LDAP
settings in reva.
    
Also we allow now to use a new set of top-level LDAP environment variables that
are shared between all LDAP-using services in ocis (graph, idp,
storage-auth-basic, storage-userprovider, storage-groupprovider, idm). This
should simplify the most LDAP based configurations considerably.

Here is a list of the new environment variables:
LDAP_URI
LDAP_INSECURE
LDAP_CACERT
LDAP_BIND_DN
LDAP_BIND_PASSWORD
LDAP_LOGIN_ATTRIBUTES
LDAP_USER_BASE_DN
LDAP_USER_SCOPE
LDAP_USER_FILTER
LDAP_USER_OBJECTCLASS
LDAP_USER_SCHEMA_MAIL
LDAP_USER_SCHEMA_DISPLAY_NAME
LDAP_USER_SCHEMA_USERNAME
LDAP_USER_SCHEMA_ID
LDAP_USER_SCHEMA_ID_IS_OCTETSTRING
LDAP_GROUP_BASE_DN
LDAP_GROUP_SCOPE
LDAP_GROUP_FILTER
LDAP_GROUP_OBJECTCLASS
LDAP_GROUP_SCHEMA_GROUPNAME
LDAP_GROUP_SCHEMA_ID
LDAP_GROUP_SCHEMA_ID_IS_OCTETSTRING

Where need these can be overwritten by service specific variables. E.g. it is possible
to use STORAGE_LDAP_URI to overide the top-level LDAP_URI variable.

https://github.com/owncloud/ocis/pull/3476
https://github.com/owncloud/ocis/issues/3150
