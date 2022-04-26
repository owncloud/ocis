## Environment Variables

| Name | Type | Default Value | Description |
|------|------|---------------|-------------|
| AUTH_BEARER_DEBUG_ADDR | string | 127.0.0.1:9149 | |
| AUTH_BEARER_DEBUG_TOKEN | string |  | |
| AUTH_BEARER_DEBUG_PPROF | bool | false | |
| AUTH_BEARER_DEBUG_ZPAGES | bool | false | |
| AUTH_BEARER_GRPC_ADDR | string | 127.0.0.1:9148 | The address of the grpc service.|
| AUTH_BEARER_GRPC_PROTOCOL | string | tcp | The transport protocol of the grpc service.|
| AUTH_BEARER_AUTH_PROVIDER | string | ldap | The auth provider which should be used by the service|
| OCIS_URL;AUTH_BEARER_OIDC_ISSUER | string | https://localhost:9200 | |
| OCIS_INSECURE;AUTH_BEARER_OIDC_INSECURE | bool | false | |