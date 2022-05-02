## Environment Variables

| Name | Type | Default Value | Description |
|------|------|---------------|-------------|
| PROXY_DEBUG_ADDR | string | 127.0.0.1:9205 | |
| PROXY_DEBUG_TOKEN | string |  | |
| PROXY_DEBUG_PPROF | bool | false | |
| PROXY_DEBUG_ZPAGES | bool | false | |
| PROXY_HTTP_ADDR | string | 0.0.0.0:9200 | |
| PROXY_HTTP_ROOT | string | / | |
| PROXY_TRANSPORT_TLS_CERT | string | ~/.ocis/proxy/server.crt | |
| PROXY_TRANSPORT_TLS_KEY | string | ~/.ocis/proxy/server.key | |
| PROXY_TLS | bool | true | |
| OCIS_URL;PROXY_OIDC_ISSUER | string | https://localhost:9200 | |
| OCIS_INSECURE;PROXY_OIDC_INSECURE | bool | true | |
| PROXY_OIDC_USERINFO_CACHE_SIZE | int | 1024 | |
| PROXY_OIDC_USERINFO_CACHE_TTL | int | 10 | |
| PROXY_ENABLE_PRESIGNEDURLS | bool | true | |
| PROXY_ACCOUNT_BACKEND_TYPE | string | cs3 | |
| PROXY_USER_OIDC_CLAIM | string | email | |
| PROXY_USER_CS3_CLAIM | string | mail | |
| OCIS_MACHINE_AUTH_API_KEY;PROXY_MACHINE_AUTH_API_KEY | string |  | |
| PROXY_AUTOPROVISION_ACCOUNTS | bool | false | |
| PROXY_ENABLE_BASIC_AUTH | bool | false | |
| PROXY_INSECURE_BACKENDS | bool | false | |