## Environment Variables

| Name | Type | Default Value | Description |
|------|------|---------------|-------------|
| OCIS_TRACING_ENABLED<br/>THUMBNAILS_TRACING_ENABLED | bool | false | Activates tracing.|
| OCIS_TRACING_TYPE<br/>THUMBNAILS_TRACING_TYPE | string |  | The type of tracing. Defaults to "", which is the same as "jaeger". Allowed tracing types are "jaeger" and "" as of now.|
| OCIS_TRACING_ENDPOINT<br/>THUMBNAILS_TRACING_ENDPOINT | string |  | The endpoint of the tracing agent.|
| OCIS_TRACING_COLLECTOR<br/>THUMBNAILS_TRACING_COLLECTOR | string |  | The HTTP endpoint for sending spans directly to a collector, i.e. http://jaeger-collector:14268/api/traces. Only used if the tracing endpoint is unset.|
| OCIS_LOG_LEVEL<br/>THUMBNAILS_LOG_LEVEL | string |  | The log level. Valid values are: "panic", "fatal", "error", "warn", "info", "debug", "trace".|
| OCIS_LOG_PRETTY<br/>THUMBNAILS_LOG_PRETTY | bool | false | Activates pretty log output.|
| OCIS_LOG_COLOR<br/>THUMBNAILS_LOG_COLOR | bool | false | Activates colorized log output.|
| OCIS_LOG_FILE<br/>THUMBNAILS_LOG_FILE | string |  | The path to the log file. Activates logging to this file if set.|
| THUMBNAILS_DEBUG_ADDR | string | 127.0.0.1:9189 | Bind address of the debug server, where metrics, health, config and debug endpoints will be exposed.|
| THUMBNAILS_DEBUG_TOKEN | string |  | Token to secure the metrics endpoint|
| THUMBNAILS_DEBUG_PPROF | bool | false | Enables pprof, which can be used for profiling|
| THUMBNAILS_DEBUG_ZPAGES | bool | false | Enables zpages, which can be used for collecting and viewing in-memory traces.|
| THUMBNAILS_GRPC_ADDR | string | 127.0.0.1:9185 | The address off the grpc service.|
| THUMBNAILS_HTTP_ADDR | string | 127.0.0.1:9186 | The bind address of the HTTP service.|
| THUMBNAILS_HTTP_ROOT | string | /thumbnails | Subdirectory that serves as the root for this HTTP service.|
| THUMBNAILS_RESOLUTIONS | []string | [16x16 32x32 64x64 128x128 1920x1080 3840x2160 7680x4320] | The supported target resolutions in the format WidthxHeight e.g. 32x32. You can define any resolution as required and separate multiple resolutions by blank or comma.|
| THUMBNAILS_FILESYSTEMSTORAGE_ROOT | string | ~/.ocis/thumbnails | The directory where the filesystem storage will store the thumbnails.|
| OCIS_INSECURE<br/>THUMBNAILS_WEBDAVSOURCE_INSECURE | bool | false | Ignore untrusted SSL certificates when connecting to the webdav source.|
| OCIS_INSECURE<br/>THUMBNAILS_CS3SOURCE_INSECURE | bool | false | Ignore untrusted SSL certificates when connecting to the CS3 source.|
| REVA_GATEWAY | string | 127.0.0.1:9142 | The CS3 gateway endpoint.|
| THUMBNAILS_TXT_FONTMAP_FILE | string |  | The path to a font file for txt thumbnails.|
| THUMBNAILS_TRANSFER_TOKEN | string |  | The secret to sign JWT to download the actual thumbnail file.|
| THUMBNAILS_DATA_ENDPOINT | string | http://127.0.0.1:9186/thumbnails/data | The HTTP endpoint where the actual thumbnail file can be downloaded.|