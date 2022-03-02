## Environment Variables

| Name | Type | Default Value | Description |
|------|------|---------------|-------------|
| IDP_DEBUG_ADDR | string | 127.0.0.1:9134 | |
| IDP_DEBUG_TOKEN | string |  | |
| IDP_DEBUG_PPROF | bool | false | |
| IDP_DEBUG_ZPAGES | bool | false | |
| IDP_HTTP_ADDR | string | 127.0.0.1:9130 | |
| IDP_HTTP_ROOT | string | / | |
| IDP_TRANSPORT_TLS_CERT | string | ~/.ocis/idp/server.crt | |
| IDP_TRANSPORT_TLS_KEY | string | ~/.ocis/idp/server.key | |
| IDP_TLS | bool | false | |
| IDP_ASSET_PATH | string |  | |
| OCIS_URL;IDP_ISS | string | https://localhost:9200 | |
| IDP_IDENTITY_MANAGER | string | ldap | |
| IDP_URI_BASE_PATH | string |  | |
| IDP_SIGN_IN_URI | string |  | |
| IDP_SIGN_OUT_URI | string |  | |
| IDP_ENDPOINT_URI | string |  | |
| IDP_ENDSESSION_ENDPOINT_URI | string |  | |
| IDP_INSECURE | bool | false | |
| IDP_ALLOW_CLIENT_GUESTS | bool | false | |
| IDP_ALLOW_DYNAMIC_CLIENT_REGISTRATION | bool | false | |
| IDP_ENCRYPTION_SECRET | string |  | |
| IDP_DISABLE_IDENTIFIER_WEBAPP | bool | true | |
| IDP_IDENTIFIER_CLIENT_PATH | string | ~/.ocis/idp | |
| IDP_IDENTIFIER_REGISTRATION_CONF | string | ~/.ocis/idp/identifier-registration.yaml | |
| IDP_IDENTIFIER_SCOPES_CONF | string |  | |
| IDP_SIGNING_KID | string |  | |
| IDP_SIGNING_METHOD | string | PS256 | |
| IDP_VALIDATION_KEYS_PATH | string |  | |
| IDP_ACCESS_TOKEN_EXPIRATION | uint64 | 600 | |
| IDP_ID_TOKEN_EXPIRATION | uint64 | 3600 | |
| IDP_REFRESH_TOKEN_EXPIRATION | uint64 | 94608000 | |
|  | uint64 | 0 | |
| IDP_LDAP_URI | string | ldap://localhost:9125 | |
| IDP_LDAP_BIND_DN | string | cn=idp,ou=sysusers,dc=ocis,dc=test | |
| IDP_LDAP_BIND_PASSWORD | string | idp | |
| IDP_LDAP_BASE_DN | string | ou=users,dc=ocis,dc=test | |
| IDP_LDAP_SCOPE | string | sub | |
| IDP_LDAP_LOGIN_ATTRIBUTE | string | cn | |
| IDP_LDAP_EMAIL_ATTRIBUTE | string | mail | |
| IDP_LDAP_NAME_ATTRIBUTE | string | sn | |
| IDP_LDAP_UUID_ATTRIBUTE | string | uid | |
| IDP_LDAP_UUID_ATTRIBUTE_TYPE | string | text | |
| IDP_LDAP_FILTER | string | (objectClass=posixaccount) | |