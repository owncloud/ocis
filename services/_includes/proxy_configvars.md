## Environment Variables

| Name | Type | Default Value | Description |
|------|------|---------------|-------------|
| OCIS_TRACING_ENABLED<br/>PROXY_TRACING_ENABLED | bool | false | Activates tracing.|
| OCIS_TRACING_TYPE<br/>PROXY_TRACING_TYPE | string |  | The type of tracing. Defaults to "", which is the same as "jaeger". Allowed tracing types are "jaeger" and "" as of now.|
| OCIS_TRACING_ENDPOINT<br/>PROXY_TRACING_ENDPOINT | string |  | The endpoint of the tracing agent.|
| OCIS_TRACING_COLLECTOR<br/>PROXY_TRACING_COLLECTOR | string |  | The HTTP endpoint for sending spans directly to a collector, i.e. http://jaeger-collector:14268/api/traces. Only used if the tracing endpoint is unset.|
| OCIS_LOG_LEVEL<br/>PROXY_LOG_LEVEL | string |  | The log level. Valid values are: "panic", "fatal", "error", "warn", "info", "debug", "trace".|
| OCIS_LOG_PRETTY<br/>PROXY_LOG_PRETTY | bool | false | Activates pretty log output.|
| OCIS_LOG_COLOR<br/>PROXY_LOG_COLOR | bool | false | Activates colorized log output.|
| OCIS_LOG_FILE<br/>PROXY_LOG_FILE | string |  | The path to the log file. Activates logging to this file if set.|
| PROXY_DEBUG_ADDR | string | 127.0.0.1:9205 | Bind address of the debug server, where metrics, health, config and debug endpoints will be exposed.|
| PROXY_DEBUG_TOKEN | string |  | Token to secure the metrics endpoint.|
| PROXY_DEBUG_PPROF | bool | false | Enables pprof, which can be used for profiling.|
| PROXY_DEBUG_ZPAGES | bool | false | Enables zpages, which can be used for collecting and viewing in-memory traces.|
| PROXY_HTTP_ADDR | string | 0.0.0.0:9200 | The bind address of the HTTP service.|
| PROXY_HTTP_ROOT | string | / | Subdirectory that serves as the root for this HTTP service.|
| PROXY_TRANSPORT_TLS_CERT | string | ~/.ocis/proxy/server.crt | File name of the TLS server certificate for the HTTPS server.|
| PROXY_TRANSPORT_TLS_KEY | string | ~/.ocis/proxy/server.key | File name of the TLS server certificate key for the HTTPS server.|
| PROXY_TLS | bool | true | Use the HTTPS server instead of the HTTP server.|
| REVA_GATEWAY | string | 127.0.0.1:9142 | The CS3 gateway endpoint.|
| OCIS_URL<br/>OCIS_OIDC_ISSUER<br/>PROXY_OIDC_ISSUER | string | https://localhost:9200 | URL of the OIDC issuer. It defaults to URL of the builtin IDP.|
| OCIS_INSECURE<br/>PROXY_OIDC_INSECURE | bool | true | Disable TLS certificate validation for connections to the IDP. Note that this is not recommended for production environments.|
| PROXY_OIDC_ACCESS_TOKEN_VERIFY_METHOD | string | jwt | Sets how OIDC access tokens should be verified. Possible values are 'none' and 'jwt'. When using 'none', no special validation apart from using it for accessing the IPD's userinfo endpoint will be done. When using 'jwt', it tries to parse the access token as a jwt token and verifies the signature using the keys published on the IDP's 'jwks_uri'.|
| PROXY_OIDC_USERINFO_CACHE_SIZE | int | 1024 | Cache size for OIDC user info.|
| PROXY_OIDC_USERINFO_CACHE_TTL | int | 10 | Max TTL in seconds for the OIDC user info cache.|
| PROXY_OIDC_JWKS_REFRESH_INTERVAL | uint64 | 60 | The interval for refreshing the JWKS (JSON Web Key Set) in minutes in the background via a new HTTP request to the IDP.|
| PROXY_OIDC_JWKS_REFRESH_TIMEOUT | uint64 | 10 | The timeout in seconds for an outgoing JWKS request.|
| PROXY_OIDC_JWKS_REFRESH_RATE_LIMIT | uint64 | 60 | Limits the rate in seconds at which refresh requests are performed for unknown keys. This is used to prevent malicious clients from imposing high network load on the IDP via ocis.|
| PROXY_OIDC_JWKS_REFRESH_UNKNOWN_KID | bool | true | If set to 'true', the JWKS refresh request will occur every time an unknown KEY ID (KID) is seen. Always set a 'refresh_limit' when enabling this.|
| OCIS_JWT_SECRET<br/>PROXY_JWT_SECRET | string |  | The secret to mint and validate JWT tokens.|
| PROXY_ENABLE_PRESIGNEDURLS | bool | true | Allow OCS to get a signing key to sign requests.|
| PROXY_ACCOUNT_BACKEND_TYPE | string | cs3 | Account backend the PROXY service should use. Currently only 'cs3' is possible here.|
| PROXY_USER_OIDC_CLAIM | string | preferred_username | The name of an OpenID Connect claim that should be used for resolving users with the account backend. Currently defaults to 'email'.|
| PROXY_USER_CS3_CLAIM | string | username | The name of a CS3 user attribute (claim) that should be mapped to the 'user_oidc_claim'. Supported values are 'username', 'mail' and 'userid'.|
| OCIS_MACHINE_AUTH_API_KEY<br/>PROXY_MACHINE_AUTH_API_KEY | string |  | Machine auth API key used to validate internal requests necessary to access resources from other services.|
| PROXY_AUTOPROVISION_ACCOUNTS | bool | false | Set this to 'true' to automatically provision users that do not yet exist in the users service on-demand upon first sign-in. To use this a write-enabled libregraph user backend needs to be setup an running.|
| PROXY_ENABLE_BASIC_AUTH | bool | false | Set this to true to enable 'basic' (username/password) authentication.|
| PROXY_INSECURE_BACKENDS | bool | false | Disable TLS certificate validation for all HTTP backend connections.|