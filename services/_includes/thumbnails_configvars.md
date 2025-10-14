## Environment Variables

| Name | Type | Default Value | Description |
|------|------|---------------|-------------|
| OCIS_TRACING_ENABLED<br/>THUMBNAILS_TRACING_ENABLED | bool | false | Activates tracing.|
| OCIS_TRACING_TYPE<br/>THUMBNAILS_TRACING_TYPE | string |  | The type of tracing. Defaults to '', which is the same as 'jaeger'. Allowed tracing types are 'jaeger', 'otlp' and '' as of now.|
| OCIS_TRACING_ENDPOINT<br/>THUMBNAILS_TRACING_ENDPOINT | string |  | The endpoint of the tracing agent.|
| OCIS_TRACING_COLLECTOR<br/>THUMBNAILS_TRACING_COLLECTOR | string |  | The HTTP endpoint for sending spans directly to a collector, i.e. http://jaeger-collector:14268/api/traces. Only used if the tracing endpoint is unset.|
| OCIS_LOG_LEVEL<br/>THUMBNAILS_LOG_LEVEL | string |  | The log level. Valid values are: 'panic', 'fatal', 'error', 'warn', 'info', 'debug', 'trace'.|
| OCIS_LOG_PRETTY<br/>THUMBNAILS_LOG_PRETTY | bool | false | Activates pretty log output.|
| OCIS_LOG_COLOR<br/>THUMBNAILS_LOG_COLOR | bool | false | Activates colorized log output.|
| OCIS_LOG_FILE<br/>THUMBNAILS_LOG_FILE | string |  | The path to the log file. Activates logging to this file if set.|
| THUMBNAILS_DEBUG_ADDR | string | 127.0.0.1:9189 | Bind address of the debug server, where metrics, health, config and debug endpoints will be exposed.|
| THUMBNAILS_DEBUG_TOKEN | string |  | Token to secure the metrics endpoint.|
| THUMBNAILS_DEBUG_PPROF | bool | false | Enables pprof, which can be used for profiling.|
| THUMBNAILS_DEBUG_ZPAGES | bool | false | Enables zpages, which can be used for collecting and viewing in-memory traces.|
| THUMBNAILS_GRPC_ADDR | string | 127.0.0.1:9185 | The bind address of the GRPC service.|
| THUMBNAILS_MAX_CONCURRENT_REQUESTS | int | 0 | Number of maximum concurrent thumbnail requests. Default is 0 which is unlimited.|
| THUMBNAILS_HTTP_ADDR | string | 127.0.0.1:9186 | The bind address of the HTTP service.|
| OCIS_HTTP_TLS_ENABLED | bool | false | Activates TLS for the http based services using the server certifcate and key configured via OCIS_HTTP_TLS_CERTIFICATE and OCIS_HTTP_TLS_KEY. If OCIS_HTTP_TLS_CERTIFICATE is not set a temporary server certificate is generated - to be used with PROXY_INSECURE_BACKEND=true.|
| OCIS_HTTP_TLS_CERTIFICATE | string |  | Path/File name of the TLS server certificate (in PEM format) for the http services.|
| OCIS_HTTP_TLS_KEY | string |  | Path/File name for the TLS certificate key (in PEM format) for the server certificate to use for the http services.|
| THUMBNAILS_HTTP_ROOT | string | /thumbnails | Subdirectory that serves as the root for this HTTP service.|
| OCIS_CORS_ALLOW_ORIGINS<br/>THUMBNAILS_CORS_ALLOW_ORIGINS | []string | [*] | A list of allowed CORS origins. See following chapter for more details: *Access-Control-Allow-Origin* at https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Origin. See the Environment Variable Types description for more details.|
| OCIS_CORS_ALLOW_METHODS<br/>THUMBNAILS_CORS_ALLOW_METHODS | []string | [GET POST PUT PATCH DELETE OPTIONS] | A list of allowed CORS methods. See following chapter for more details: *Access-Control-Request-Method* at https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Request-Method. See the Environment Variable Types description for more details.|
| OCIS_CORS_ALLOW_HEADERS<br/>THUMBNAILS_CORS_ALLOW_HEADERS | []string | [Authorization Origin Content-Type Accept X-Requested-With X-Request-Id Cache-Control] | A list of allowed CORS headers. See following chapter for more details: *Access-Control-Request-Headers* at https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Request-Headers. See the Environment Variable Types description for more details.|
| OCIS_CORS_ALLOW_CREDENTIALS<br/>THUMBNAILS_CORS_ALLOW_CREDENTIALS | bool | true | Allow credentials for CORS.See following chapter for more details: *Access-Control-Allow-Credentials* at https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Credentials.|
| THUMBNAILS_RESOLUTIONS | []string | [16x16 32x32 64x64 128x128 1080x1920 1920x1080 2160x3840 3840x2160 4320x7680 7680x4320] | The supported list of target resolutions in the format WidthxHeight like 32x32. You can define any resolution as required. See the Environment Variable Types description for more details.|
| THUMBNAILS_FILESYSTEMSTORAGE_ROOT | string | /var/lib/ocis/thumbnails | The directory where the filesystem storage will store the thumbnails. If not defined, the root directory derives from $OCIS_BASE_DATA_PATH/thumbnails.|
| OCIS_INSECURE<br/>THUMBNAILS_WEBDAVSOURCE_INSECURE | bool | false | Ignore untrusted SSL certificates when connecting to the webdav source.|
| OCIS_INSECURE<br/>THUMBNAILS_CS3SOURCE_INSECURE | bool | false | Ignore untrusted SSL certificates when connecting to the CS3 source.|
| OCIS_REVA_GATEWAY | string | com.owncloud.api.gateway | CS3 gateway used to look up user metadata|
| THUMBNAILS_TXT_FONTMAP_FILE | string |  | The path to a font file for txt thumbnails.|
| THUMBNAILS_TRANSFER_TOKEN | string |  | The secret to sign JWT to download the actual thumbnail file.|
| THUMBNAILS_DATA_ENDPOINT | string | http://127.0.0.1:9186/thumbnails/data | The HTTP endpoint where the actual thumbnail file can be downloaded.|
| THUMBNAILS_MAX_INPUT_WIDTH | int | 7680 | The maximum width of an input image which is being processed.|
| THUMBNAILS_MAX_INPUT_HEIGHT | int | 7680 | The maximum height of an input image which is being processed.|
| THUMBNAILS_MAX_INPUT_IMAGE_FILE_SIZE | string | 50MB | The maximum file size of an input image which is being processed. Usable common abbreviations: [KB, KiB, MB, MiB, GB, GiB, TB, TiB, PB, PiB, EB, EiB], example: 2GB.|