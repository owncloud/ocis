## Environment Variables

| Name | Type | Default Value | Description |
|------|------|---------------|-------------|
| IDM_DEBUG_ADDR | string |  | |
| IDM_DEBUG_TOKEN | string |  | |
| IDM_DEBUG_PPROF | bool | false | |
| IDM_DEBUG_ZPAGES | bool | false | |
| IDM_LDAPS_ADDR | string | 127.0.0.1:9235 | Listen address for the ldaps listener (ip-addr:port)|
| IDM_LDAPS_CERT | string | ~/.ocis/idm/ldap.crt | File name of the TLS server certificate for the ldaps listener|
| IDM_LDAPS_KEY | string | ~/.ocis/idm/ldap.key | File name for the TLS certificate key for the server certificate|
| IDM_DATABASE_PATH | string | ~/.ocis/idm/ocis.boltdb | Full path to the idm backend database|
| IDM_CREATE_DEMO_USERS;ACCOUNTS_DEMO_USERS_AND_GROUPS | bool | true | Flag to enabe/disable the creation of the demo users|
| IDM_ADMIN_PASSWORD | string | idm | Password to set for the "idm" service users. Either cleartext or an argon2id hash|
| IDM_REVASVC_PASSWORD | string | reva | Password to set for the "reva" service users. Either cleartext or an argon2id hash|
| IDM_IDPSVC_PASSWORD | string | idp | Password to set for the "idp" service users. Either cleartext or an argon2id hash|