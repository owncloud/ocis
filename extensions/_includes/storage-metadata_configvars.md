## Environment Variables

| Name | Type | Default Value | Description |
|------|------|---------------|-------------|
| STORAGE_METADATA_DEBUG_ADDR | string | 127.0.0.1:9217 | |
| STORAGE_METADATA_DEBUG_TOKEN | string |  | |
| STORAGE_METADATA_DEBUG_PPROF | bool | false | |
| STORAGE_METADATA_DEBUG_ZPAGES | bool | false | |
| STORAGE_METADATA_GRPC_ADDR | string | 127.0.0.1:9215 | The address of the grpc service.|
| STORAGE_METADATA_GRPC_PROTOCOL | string | tcp | The transport protocol of the grpc service.|
| STORAGE_METADATA_GRPC_ADDR | string | 127.0.0.1:9216 | The address of the grpc service.|
| STORAGE_METADATA_GRPC_PROTOCOL | string | tcp | The transport protocol of the grpc service.|
| STORAGE_METADATA_DRIVER | string | ocis | The driver which should be used by the service|
| STORAGE_METADATA_DRIVER_OCIS_ROOT | string | ~/.ocis/storage/metadata | |
| OCIS_INSECURE;STORAGE_METADATA_DATAPROVIDER_INSECURE | bool | false | |