## Environment Variables

| Name | Type | Default Value | Description |
|------|------|---------------|-------------|
| COLLABORATION_APP_NAME | string | WOPI app | The name of the app|
| COLLABORATION_APP_DESCRIPTION | string | Open office documents with a WOPI app | App description|
| COLLABORATION_APP_ICON | string | image-edit | Icon for the app|
| COLLABORATION_APP_LOCKNAME | string | com.github.owncloud.collaboration | Name for the app lock|
| OCIS_JWT_SECRET<br/>COLLABORATION_JWT_SECRET | string | oRy1VFYC1w3Y0LWoUp2PstWFYgcc4Uem | Used as JWT token and to encrypt access token.|
| COLLABORATION_GRPC_ADDR | string | 0.0.0.0:9301 | The bind address of the GRPC service.|
| COLLABORATION_HTTP_ADDR | string | 127.0.0.1:9300 | The external address of the collaboration service wihout a leading scheme. Either use an IP address or a hostname (127.0.0.1:9301 or wopi.private.prv). The configured 'Scheme' in another envvar will be used to finally build the public URL along with this address.|
| COLLABORATION_HTTP_BINDADDR | string | 0.0.0.0:9300 | The bind address of the HTTP service. Use '<ip-address>:<port>', for example, '127.0.0.1:9301' or '0.0.0.0:9301'.|
| COLLABORATION_HTTP_SCHEME | string | https | The scheme to use for the HTTP address, which is either 'http' or 'https'.|
| OCIS_HTTP_TLS_ENABLED | bool | false | Activates TLS for the http based services using the server certifcate and key configured via OCIS_HTTP_TLS_CERTIFICATE and OCIS_HTTP_TLS_KEY. If OCIS_HTTP_TLS_CERTIFICATE is not set a temporary server certificate is generated - to be used with PROXY_INSECURE_BACKEND=true.|
| OCIS_HTTP_TLS_CERTIFICATE | string |  | Path/File name of the TLS server certificate (in PEM format) for the http services.|
| OCIS_HTTP_TLS_KEY | string |  | Path/File name for the TLS certificate key (in PEM format) for the server certificate to use for the http services.|
| COLLABORATION_WOPIAPP_ADDR | string | https://127.0.0.1:8080 | The URL where the WOPI app is located, such as https://127.0.0.1:8080.|
| COLLABORATION_WOPIAPP_INSECURE | bool | false | Skip TLS certificate verification when connecting to the WOPI app|
| COLLABORATION_CS3API_DATAGATEWAY_INSECURE | bool | false | Connect to the CS3API data gateway insecurely.|
| OCIS_TRACING_ENABLED<br/>COLLABORATION_TRACING_ENABLED | bool | false | Activates tracing.|
| OCIS_TRACING_TYPE<br/>COLLABORATION_TRACING_TYPE | string |  | The type of tracing. Defaults to '', which is the same as 'jaeger'. Allowed tracing types are 'jaeger' and '' as of now.|
| OCIS_TRACING_ENDPOINT<br/>COLLABORATION_TRACING_ENDPOINT | string |  | The endpoint of the tracing agent.|
| OCIS_TRACING_COLLECTOR<br/>COLLABORATION_TRACING_COLLECTOR | string |  | The HTTP endpoint for sending spans directly to a collector, i.e. http://jaeger-collector:14268/api/traces. Only used if the tracing endpoint is unset.|
| OCIS_LOG_LEVEL<br/>COLLABORATION_LOG_LEVEL | string |  | The log level. Valid values are: 'panic', 'fatal', 'error', 'warn', 'info', 'debug', 'trace'.|
| OCIS_LOG_PRETTY<br/>COLLABORATION_LOG_PRETTY | bool | false | Activates pretty log output.|
| OCIS_LOG_COLOR<br/>COLLABORATION_LOG_COLOR | bool | false | Activates colorized log output.|
| OCIS_LOG_FILE<br/>COLLABORATION_LOG_FILE | string |  | The path to the log file. Activates logging to this file if set.|
| COLLABORATION_DEBUG_ADDR | string | 127.0.0.1:9304 | Bind address of the debug server, where metrics, health, config and debug endpoints will be exposed.|
| COLLABORATION_DEBUG_TOKEN | string |  | Token to secure the metrics endpoint.|
| COLLABORATION_DEBUG_PPROF | bool | false | Enables pprof, which can be used for profiling.|
| COLLABORATION_DEBUG_ZPAGES | bool | false | Enables zpages, which can be used for collecting and viewing in-memory traces.|