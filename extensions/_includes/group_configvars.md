## Environment Variables

| Name | Type | Default Value | Description |
|------|------|---------------|-------------|
| GROUPS_DEBUG_ADDR | string | 127.0.0.1:9161 | |
| GROUPS_DEBUG_TOKEN | string |  | |
| GROUPS_DEBUG_PPROF | bool | false | |
| GROUPS_DEBUG_ZPAGES | bool | false | |
| GROUPS_GRPC_ADDR | string | 127.0.0.1:9160 | The address of the grpc service.|
| GROUPS_GRPC_PROTOCOL | string | tcp | The transport protocol of the grpc service.|
| LDAP_URI;GROUPS_LDAP_URI | string | ldaps://localhost:9235 | |
| LDAP_CACERT;GROUPS_LDAP_CACERT | string | ~/.ocis/idm/ldap.crt | |
| LDAP_INSECURE;GROUPS_LDAP_INSECURE | bool | false | |
| LDAP_BIND_DN;GROUPS_LDAP_BIND_DN | string | uid=reva,ou=sysusers,o=libregraph-idm | |
| LDAP_BIND_PASSWORD;GROUPS_LDAP_BIND_PASSWORD | string |  | |
| LDAP_USER_BASE_DN;GROUPS_LDAP_USER_BASE_DN | string | ou=users,o=libregraph-idm | |
| LDAP_GROUP_BASE_DN;GROUPS_LDAP_GROUP_BASE_DN | string | ou=groups,o=libregraph-idm | |
| LDAP_USER_SCOPE;GROUPS_LDAP_USER_SCOPE | string | sub | |
| LDAP_GROUP_SCOPE;GROUPS_LDAP_GROUP_SCOPE | string | sub | |
| LDAP_USERFILTER;GROUPS_LDAP_USERFILTER | string |  | |
| LDAP_GROUPFILTER;GROUPS_LDAP_USERFILTER | string |  | |
| LDAP_USER_OBJECTCLASS;GROUPS_LDAP_USER_OBJECTCLASS | string | inetOrgPerson | |
| LDAP_GROUP_OBJECTCLASS;GROUPS_LDAP_GROUP_OBJECTCLASS | string | groupOfNames | |
| LDAP_LOGIN_ATTRIBUTES;GROUPS_LDAP_LOGIN_ATTRIBUTES |  | [uid mail] | |
| OCIS_URL;GROUPS_IDP_URL | string | https://localhost:9200 | |
| LDAP_USER_SCHEMA_ID;GROUPS_LDAP_USER_SCHEMA_ID | string | ownclouduuid | |
| LDAP_USER_SCHEMA_ID_IS_OCTETSTRING;GROUPS_LDAP_USER_SCHEMA_ID_IS_OCTETSTRING | bool | false | |
| LDAP_USER_SCHEMA_MAIL;GROUPS_LDAP_USER_SCHEMA_MAIL | string | mail | |
| LDAP_USER_SCHEMA_DISPLAYNAME;GROUPS_LDAP_USER_SCHEMA_DISPLAYNAME | string | displayname | |
| LDAP_USER_SCHEMA_USERNAME;GROUPS_LDAP_USER_SCHEMA_USERNAME | string | uid | |
| LDAP_GROUP_SCHEMA_ID;GROUPS_LDAP_GROUP_SCHEMA_ID | string | ownclouduuid | |
| LDAP_GROUP_SCHEMA_ID_IS_OCTETSTRING;GROUPS_LDAP_GROUP_SCHEMA_ID_IS_OCTETSTRING | bool | false | |
| LDAP_GROUP_SCHEMA_MAIL;GROUPS_LDAP_GROUP_SCHEMA_MAIL | string | mail | |
| LDAP_GROUP_SCHEMA_DISPLAYNAME;GROUPS_LDAP_GROUP_SCHEMA_DISPLAYNAME | string | cn | |
| LDAP_GROUP_SCHEMA_GROUPNAME;GROUPS_LDAP_GROUP_SCHEMA_GROUPNAME | string | cn | |
| LDAP_GROUP_SCHEMA_MEMBER;GROUPS_LDAP_GROUP_SCHEMA_MEMBER | string | member | |