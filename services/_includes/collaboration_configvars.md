## Environment Variables

| Name | Type | Default Value | Description |
|------|------|---------------|-------------|
| COLLABORATION_APP_NAME | string | Collabora | The name of the app which is shown to the user. You can chose freely but you are limited to a single word without special characters or whitespaces. We recommend to use pascalCase like 'CollaboraOnline'.|
| COLLABORATION_APP_PRODUCT | string | Collabora | The WebOffice app, either Collabora, OnlyOffice, Microsoft365 or MicrosoftOfficeOnline.|
| COLLABORATION_APP_PRODUCT_EDITION | string |  | The WebOffice app edition defines the capabilities specific to the product such as CE - Community Edition, EE - Enterprise Edition DE - Developer Edition, etc. Currently supported values are limited to OnlyOffice and are: 'ce', 'ee' or 'de' and  default to empty which is equal to ce). See the documentation for more details.|
| COLLABORATION_APP_DESCRIPTION | string | Open office documents with Collabora | App description|
| COLLABORATION_APP_ICON | string | image-edit | Icon for the app|
| COLLABORATION_APP_ADDR | string | https://127.0.0.1:9980 | The URL where the WOPI app is located, such as https://127.0.0.1:8080.|
| COLLABORATION_APP_INSECURE | bool | false | Skip TLS certificate verification when connecting to the WOPI app|
| COLLABORATION_APP_PROOF_DISABLE | bool | false | Disable the proof keys verification|
| COLLABORATION_APP_PROOF_DURATION | string | 12h | Duration for the proof keys to be cached in memory, using time.ParseDuration format. If the duration can't be parsed, we'll use the default 12h as duration|
| COLLABORATION_APP_LICENSE_CHECK_ENABLE | bool | false | Enable license checking to edit files. Needs to be enabled when using Microsoft365 with the business flow.|
| OCIS_PERSISTENT_STORE<br/>COLLABORATION_STORE | string | nats-js-kv | The type of the store. Supported values are: 'memory', 'nats-js-kv', 'redis-sentinel', 'noop'. See the text description for details.|
| OCIS_PERSISTENT_STORE_NODES<br/>COLLABORATION_STORE_NODES | []string | [127.0.0.1:9233] | A list of nodes to access the configured store. This has no effect when 'memory' store is configured. Note that the behaviour how nodes are used is dependent on the library of the configured store. See the Environment Variable Types description for more details.|
| COLLABORATION_STORE_DATABASE | string | collaboration | The database name the configured store should use.|
| COLLABORATION_STORE_TABLE | string |  | The database table the store should use.|
| OCIS_PERSISTENT_STORE_TTL<br/>COLLABORATION_STORE_TTL | Duration | 30m0s | Time to live for events in the store. Defaults to '30m' (30 minutes). See the Environment Variable Types description for more details.|
| OCIS_PERSISTENT_STORE_AUTH_USERNAME<br/>COLLABORATION_STORE_AUTH_USERNAME | string |  | The username to authenticate with the store. Only applies when store type 'nats-js-kv' is configured.|
| OCIS_PERSISTENT_STORE_AUTH_PASSWORD<br/>COLLABORATION_STORE_AUTH_PASSWORD | string |  | The password to authenticate with the store. Only applies when store type 'nats-js-kv' is configured.|
| OCIS_JWT_SECRET<br/>COLLABORATION_JWT_SECRET | string |  | The secret to mint and validate jwt tokens.|
| COLLABORATION_GRPC_ADDR | string | 127.0.0.1:9301 | The bind address of the GRPC service.|
| OCIS_GRPC_PROTOCOL<br/>COLLABORATION_GRPC_PROTOCOL | string | tcp | The transport protocol of the GRPC service.|
| COLLABORATION_HTTP_ADDR | string | 127.0.0.1:9300 | The bind address of the HTTP service.|
| OCIS_HTTP_TLS_ENABLED | bool | false | Activates TLS for the http based services using the server certifcate and key configured via OCIS_HTTP_TLS_CERTIFICATE and OCIS_HTTP_TLS_KEY. If OCIS_HTTP_TLS_CERTIFICATE is not set a temporary server certificate is generated - to be used with PROXY_INSECURE_BACKEND=true.|
| OCIS_HTTP_TLS_CERTIFICATE | string |  | Path/File name of the TLS server certificate (in PEM format) for the http services.|
| OCIS_HTTP_TLS_KEY | string |  | Path/File name for the TLS certificate key (in PEM format) for the server certificate to use for the http services.|
| COLLABORATION_WOPI_SRC | string | https://localhost:9300 | The WOPI source base URL containing schema, host and port. Set this to the schema and domain where the collaboration service is reachable for the wopi app, such as https://office.owncloud.test.|
| COLLABORATION_WOPI_SECRET | string |  | Used to mint and verify WOPI JWT tokens and encrypt and decrypt the REVA JWT token embedded in the WOPI JWT token.|
| COLLABORATION_WOPI_DISABLE_CHAT<br/>OCIS_WOPI_DISABLE_CHAT | bool | false | Disable chat in the office web frontend. This feature applies to OnlyOffice and Microsoft.|
| COLLABORATION_WOPI_ENABLE_MOBILE | bool | false | Enable the mobile web view for office app. This feature applies to OnlyOffice.  See the documentation for more details.|
| COLLABORATION_WOPI_PROXY_URL | string |  | The URL to the ownCloud Office365 WOPI proxy. Optional. To use this feature, you need an office365 proxy subscription. If you become part of the Microsoft CSP program (https://learn.microsoft.com/en-us/partner-center/enroll/csp-overview), you can use WebOffice without a proxy.|
| COLLABORATION_WOPI_PROXY_SECRET | string |  | Optional, the secret to authenticate against the ownCloud Office365 WOPI proxy. This secret can be obtained from ownCloud via the office365 proxy subscription.|
| COLLABORATION_WOPI_SHORTTOKENS | bool | false | Use short access tokens for WOPI access. This is useful for office packages, like Microsoft Office Online, which have URL length restrictions. If enabled, a persistent store must be configured.|
| COLLABORATION_WOPI_DISABLED_EXTENSIONS | []string | [] | List of extensions to disable: Disabling an extension will make it unavailable to the Office web front end. The list is comma-separated with no spaces between the items, such as 'docx,xlsx,pptx'.|
| OCIS_REVA_GATEWAY | string | com.owncloud.api.gateway | CS3 gateway used to look up user metadata.|
| COLLABORATION_CS3API_DATAGATEWAY_INSECURE | bool | false | Connect to the CS3API data gateway insecurely.|
| OCIS_TRACING_ENABLED<br/>COLLABORATION_TRACING_ENABLED | bool | false | Activates tracing.|
| OCIS_TRACING_TYPE<br/>COLLABORATION_TRACING_TYPE | string |  | The type of tracing. Defaults to '', which is the same as 'jaeger'. Allowed tracing types are 'jaeger', 'otlp' and '' as of now.|
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