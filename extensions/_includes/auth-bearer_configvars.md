## Environment Variables

| Name | Type | Default Value | Description |
|------|------|---------------|-------------|
| OCIS_TRACING_ENABLED<br/>AUTH_BEARER_TRACING_ENABLED | bool | false | Activates tracing.|
| OCIS_TRACING_TYPE<br/>AUTH_BEARER_TRACING_TYPE | string |  | |
| OCIS_TRACING_ENDPOINT<br/>AUTH_BEARER_TRACING_ENDPOINT | string |  | The endpoint to the tracing collector.|
| OCIS_TRACING_COLLECTOR<br/>AUTH_BEARER_TRACING_COLLECTOR | string |  | |
| OCIS_LOG_LEVEL<br/>AUTH_BEARER_LOG_LEVEL | string |  | The log level.|
| OCIS_LOG_PRETTY<br/>AUTH_BEARER_LOG_PRETTY | bool | false | Activates pretty log output.|
| OCIS_LOG_COLOR<br/>AUTH_BEARER_LOG_COLOR | bool | false | Activates colorized log output.|
| OCIS_LOG_FILE<br/>AUTH_BEARER_LOG_FILE | string |  | The target log file.|
| AUTH_BEARER_DEBUG_ADDR | string | 127.0.0.1:9149 | |
| AUTH_BEARER_DEBUG_TOKEN | string |  | |
| AUTH_BEARER_DEBUG_PPROF | bool | false | |
| AUTH_BEARER_DEBUG_ZPAGES | bool | false | |
| AUTH_BEARER_GRPC_ADDR | string | 127.0.0.1:9148 | The address of the grpc service.|
| AUTH_BEARER_GRPC_PROTOCOL | string | tcp | The transport protocol of the grpc service.|
| OCIS_JWT_SECRET<br/>AUTH_BEARER_JWT_SECRET | string |  | |
| REVA_GATEWAY | string | 127.0.0.1:9142 | |
| AUTH_BEARER_AUTH_PROVIDER | string | ldap | The auth provider which should be used by the service|
| OCIS_URL<br/>AUTH_BEARER_OIDC_ISSUER | string | https://localhost:9200 | |
| OCIS_INSECURE<br/>AUTH_BEARER_OIDC_INSECURE | bool | false | |