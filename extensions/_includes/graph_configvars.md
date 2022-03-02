## Environment Variables

| Name | Type | Default Value | Description |
|------|------|---------------|-------------|
| GRAPH_DEBUG_ADDR | string | 127.0.0.1:9124 | |
| GRAPH_DEBUG_TOKEN | string |  | |
| GRAPH_DEBUG_PPROF | bool | false | |
| GRAPH_DEBUG_ZPAGES | bool | false | |
| GRAPH_HTTP_ADDR | string | 127.0.0.1:9120 | |
| GRAPH_HTTP_ROOT | string | /graph | |
| REVA_GATEWAY | string | 127.0.0.1:9142 | |
| OCIS_JWT_SECRET;OCS_JWT_SECRET | string | Pive-Fumkiu4 | |
| OCIS_URL;GRAPH_SPACES_WEBDAV_BASE | string | https://localhost:9200 | |
| GRAPH_SPACES_WEBDAV_PATH | string | /dav/spaces/ | |
| GRAPH_SPACES_DEFAULT_QUOTA | string | 1000000000 | |
| OCIS_INSECURE;GRAPH_SPACES_INSECURE | bool | false | |
| GRAPH_SPACES_EXTENDED_SPACE_PROPERTIES_CACHE_TTL | int | 0 | |
| GRAPH_IDENTITY_BACKEND | string | cs3 | |
| GRAPH_LDAP_URI | string | ldap://localhost:9125 | |
| OCIS_INSECURE;GRAPH_LDAP_INSECURE | bool | false | |
| GRAPH_LDAP_BIND_DN | string |  | |
| GRAPH_LDAP_BIND_PASSWORD | string |  | |
| GRAPH_LDAP_SERVER_UUID | bool | false | |
| GRAPH_LDAP_SERVER_WRITE_ENABLED | bool | false | |
| GRAPH_LDAP_USER_BASE_DN | string | ou=users,dc=ocis,dc=test | |
| GRAPH_LDAP_USER_SCOPE | string | sub | |
| GRAPH_LDAP_USER_FILTER | string | (objectClass=inetOrgPerson) | |
| GRAPH_LDAP_USER_EMAIL_ATTRIBUTE | string | mail | |
| GRAPH_LDAP_USER_DISPLAYNAME_ATTRIBUTE | string | displayName | |
| GRAPH_LDAP_USER_NAME_ATTRIBUTE | string | uid | |
| GRAPH_LDAP_USER_UID_ATTRIBUTE | string | owncloudUUID | |
| GRAPH_LDAP_GROUP_BASE_DN | string | ou=groups,dc=ocis,dc=test | |
| GRAPH_LDAP_GROUP_SEARCH_SCOPE | string | sub | |
| GRAPH_LDAP_GROUP_FILTER | string | (objectclass=groupOfNames) | |
| GRAPH_LDAP_GROUP_NAME_ATTRIBUTE | string | cn | |
| GRAPH_LDAP_GROUP_ID_ATTRIBUTE | string | owncloudUUID | |