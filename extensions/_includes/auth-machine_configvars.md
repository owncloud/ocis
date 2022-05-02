## Environment Variables

| Name | Type | Default Value | Description |
|------|------|---------------|-------------|
| AUTH_MACHINE_DEBUG_ADDR | string | 127.0.0.1:9167 | |
| AUTH_MACHINE_DEBUG_TOKEN | string |  | |
| AUTH_MACHINE_DEBUG_PPROF | bool | false | |
| AUTH_MACHINE_DEBUG_ZPAGES | bool | false | |
| AUTH_MACHINE_GRPC_ADDR | string | 127.0.0.1:9166 | The address of the grpc service.|
| AUTH_MACHINE_GRPC_PROTOCOL | string | tcp | The transport protocol of the grpc service.|
| AUTH_MACHINE_AUTH_PROVIDER | string | ldap | The auth provider which should be used by the service|
| OCIS_MACHINE_AUTH_API_KEY;AUTH_MACHINE_PROVIDER_API_KEY | string |  | The api key for the machine auth provider.|