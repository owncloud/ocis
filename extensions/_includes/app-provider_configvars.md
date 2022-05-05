## Environment Variables

| Name | Type | Default Value | Description |
|------|------|---------------|-------------|
| OCIS_TRACING_ENABLED<br/>APP_PROVIDER_TRACING_ENABLED | bool | false | Activates tracing.|
| OCIS_TRACING_TYPE<br/>APP_PROVIDER_TRACING_TYPE | string |  | |
| OCIS_TRACING_ENDPOINT<br/>APP_PROVIDER_TRACING_ENDPOINT | string |  | The endpoint to the tracing collector.|
| OCIS_TRACING_COLLECTOR<br/>APP_PROVIDER_TRACING_COLLECTOR | string |  | |
| OCIS_LOG_LEVEL<br/>APP_PROVIDER_LOG_LEVEL | string |  | The log level.|
| OCIS_LOG_PRETTY<br/>APP_PROVIDER_LOG_PRETTY | bool | false | Activates pretty log output.|
| OCIS_LOG_COLOR<br/>APP_PROVIDER_LOG_COLOR | bool | false | Activates colorized log output.|
| OCIS_LOG_FILE<br/>APP_PROVIDER_LOG_FILE | string |  | The target log file.|
| APP_PROVIDER_DEBUG_ADDR | string | 127.0.0.1:9165 | |
| APP_PROVIDER_DEBUG_TOKEN | string |  | |
| APP_PROVIDER_DEBUG_PPROF | bool | false | |
| APP_PROVIDER_DEBUG_ZPAGES | bool | false | |
| APP_PROVIDER_GRPC_ADDR | string | 127.0.0.1:9164 | The address of the grpc service.|
| APP_PROVIDER_GRPC_PROTOCOL | string | tcp | The transport protocol of the grpc service.|
| OCIS_JWT_SECRET<br/>APP_PROVIDER_JWT_SECRET | string |  | |
| REVA_GATEWAY | string | 127.0.0.1:9142 | |
| APP_PROVIDER_EXTERNAL_ADDR | string |  | |
| APP_PROVIDER_DRIVER | string |  | |
| APP_PROVIDER_WOPI_APP_API_KEY | string |  | api key for the wopi app|
| APP_PROVIDER_WOPI_APP_DESKTOP_ONLY | bool | false | offer this app only on desktop|
| APP_PROVIDER_WOPI_APP_ICON_URI | string |  | uri to an app icon to be used by clients|
| APP_PROVIDER_WOPI_APP_INTERNAL_URL | string |  | internal url to the app, eg in your DMZ|
| APP_PROVIDER_WOPI_APP_NAME | string |  | human readable app name|
| APP_PROVIDER_WOPI_APP_URL | string |  | url for end users to access the app|
| APP_PROVIDER_WOPI_INSECURE | bool | false | allow insecure connections to the app|
| APP_PROVIDER_WOPI_WOPI_SERVER_IOP_SECRET | string |  | shared secret of the CS3org WOPI server|
| APP_PROVIDER_WOPI_WOPI_SERVER_EXTERNAL_URL | string |  | external url of the CS3org WOPI server|