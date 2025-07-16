## Environment Variables

| Name | Type | Default Value | Description |
|------|------|---------------|-------------|
| OCIS_TRACING_ENABLED<br/>OCM_TRACING_ENABLED | bool | false | Activates tracing.|
| OCIS_TRACING_TYPE<br/>OCM_TRACING_TYPE | string |  | The type of tracing. Defaults to '', which is the same as 'jaeger'. Allowed tracing types are 'jaeger', 'otlp' and '' as of now.|
| OCIS_TRACING_ENDPOINT<br/>OCM_TRACING_ENDPOINT | string |  | The endpoint of the tracing agent.|
| OCIS_TRACING_COLLECTOR<br/>OCM_TRACING_COLLECTOR | string |  | The HTTP endpoint for sending spans directly to a collector, i.e. http://jaeger-collector:14268/api/traces. Only used if the tracing endpoint is unset.|
| OCIS_LOG_LEVEL<br/>OCM_LOG_LEVEL | string |  | The log level. Valid values are: 'panic', 'fatal', 'error', 'warn', 'info', 'debug', 'trace'.|
| OCIS_LOG_PRETTY<br/>OCM_LOG_PRETTY | bool | false | Activates pretty log output.|
| OCIS_LOG_COLOR<br/>OCM_LOG_COLOR | bool | false | Activates colorized log output.|
| OCIS_LOG_FILE<br/>OCM_LOG_FILE | string |  | The path to the log file. Activates logging to this file if set.|
| OCM_DEBUG_ADDR | string | 127.0.0.1:9281 | Bind address of the debug server, where metrics, health, config and debug endpoints will be exposed.|
| OCM_DEBUG_TOKEN | string |  | Token to secure the metrics endpoint.|
| OCM_DEBUG_PPROF | bool | false | Enables pprof, which can be used for profiling.|
| OCM_DEBUG_ZPAGES | bool | false | Enables zpages, which can be used for collecting and viewing in-memory traces.|
| OCM_HTTP_ADDR | string | 127.0.0.1:9280 | The bind address of the HTTP service.|
| OCM_HTTP_PROTOCOL | string | tcp | The transport protocol of the HTTP service.|
| OCM_HTTP_PREFIX | string |  | The path prefix where OCM can be accessed (defaults to /).|
| OCIS_CORS_ALLOW_ORIGINS<br/>OCM_CORS_ALLOW_ORIGINS | []string | [https://localhost:9200] | A list of allowed CORS origins. See following chapter for more details: *Access-Control-Allow-Origin* at https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Origin. See the Environment Variable Types description for more details.|
| OCIS_CORS_ALLOW_METHODS<br/>OCM_CORS_ALLOW_METHODS | []string | [OPTIONS HEAD GET PUT POST DELETE MKCOL PROPFIND PROPPATCH MOVE COPY REPORT SEARCH] | A list of allowed CORS methods. See following chapter for more details: *Access-Control-Request-Method* at https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Request-Method. See the Environment Variable Types description for more details.|
| OCIS_CORS_ALLOW_HEADERS<br/>OCM_CORS_ALLOW_HEADERS | []string | [Origin Accept Content-Type Depth Authorization Ocs-Apirequest If-None-Match If-Match Destination Overwrite X-Request-Id X-Requested-With Tus-Resumable Tus-Checksum-Algorithm Upload-Concat Upload-Length Upload-Metadata Upload-Defer-Length Upload-Expires Upload-Checksum Upload-Offset X-HTTP-Method-Override Cache-Control] | A list of allowed CORS headers. See following chapter for more details: *Access-Control-Request-Headers* at https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Request-Headers. See the Environment Variable Types description for more details.|
| OCIS_CORS_ALLOW_CREDENTIALS<br/>OCM_CORS_ALLOW_CREDENTIALS | bool | false | Allow credentials for CORS.See following chapter for more details: *Access-Control-Allow-Credentials* at https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Credentials.|
| OCM_GRPC_ADDR | string | 127.0.0.1:9282 | The bind address of the GRPC service.|
| OCIS_GRPC_PROTOCOL<br/>OCM_GRPC_PROTOCOL | string |  | The transport protocol of the GRPC service.|
| OCIS_SERVICE_ACCOUNT_ID<br/>OCM_SERVICE_ACCOUNT_ID | string |  | The ID of the service account the service should use. See the 'auth-service' service description for more details.|
| OCIS_SERVICE_ACCOUNT_SECRET<br/>OCM_SERVICE_ACCOUNT_SECRET | string |  | The service account secret.|
| OCIS_EVENTS_ENDPOINT<br/>OCM_EVENTS_ENDPOINT | string | 127.0.0.1:9233 | The address of the event system. The event system is the message queuing service. It is used as message broker for the microservice architecture.|
| OCIS_EVENTS_CLUSTER<br/>OCM_EVENTS_CLUSTER | string | ocis-cluster | The clusterID of the event system. The event system is the message queuing service. It is used as message broker for the microservice architecture. Mandatory when using NATS as event system.|
| OCIS_INSECURE<br/>OCM_EVENTS_TLS_INSECURE | bool | false | Whether to verify the server TLS certificates.|
| OCIS_EVENTS_TLS_ROOT_CA_CERTIFICATE<br/>OCM_EVENTS_TLS_ROOT_CA_CERTIFICATE | string |  | The root CA certificate used to validate the server's TLS certificate. If provided OCM_EVENTS_TLS_INSECURE will be seen as false.|
| OCIS_EVENTS_ENABLE_TLS<br/>OCM_EVENTS_ENABLE_TLS | bool | false | Enable TLS for the connection to the events broker. The events broker is the ocis service which receives and delivers events between the services.|
| OCIS_EVENTS_AUTH_USERNAME<br/>OCM_EVENTS_AUTH_USERNAME | string |  | The username to authenticate with the events broker. The events broker is the ocis service which receives and delivers events between the services.|
| OCIS_EVENTS_AUTH_PASSWORD<br/>OCM_EVENTS_AUTH_PASSWORD | string |  | The password to authenticate with the events broker. The events broker is the ocis service which receives and delivers events between the services.|
| OCIS_JWT_SECRET<br/>OCM_JWT_SECRET | string |  | The secret to mint and validate jwt tokens.|
| OCIS_REVA_GATEWAY | string | com.owncloud.api.gateway | The CS3 gateway endpoint.|
| OCIS_GRPC_CLIENT_TLS_MODE | string |  | TLS mode for grpc connection to the go-micro based grpc services. Possible values are 'off', 'insecure' and 'on'. 'off': disables transport security for the clients. 'insecure' allows using transport security, but disables certificate verification (to be used with the autogenerated self-signed certificates). 'on' enables transport security, including server certificate verification.|
| OCIS_GRPC_CLIENT_TLS_CACERT | string |  | Path/File name for the root CA certificate (in PEM format) used to validate TLS server certificates of the go-micro based grpc services.|
| OCM_OCMD_PREFIX | string | ocm | URL path prefix for the OCMD service. Note that the string must not start with '/'.|
| OCM_OCMD_EXPOSE_RECIPIENT_DISPLAY_NAME | bool | false | Expose the display name of OCM share recipients.|
| OCM_SCIENCEMESH_PREFIX | string | sciencemesh | URL path prefix for the ScienceMesh service. Note that the string must not start with '/'.|
| OCM_MESH_DIRECTORY_URL | string |  | URL of the mesh directory service.|
| OCM_OCM_INVITE_MANAGER_DRIVER | string | json | Driver to be used to persist OCM invites. Supported value is only 'json'.|
| OCM_OCM_INVITE_MANAGER_JSON_FILE | string | /var/lib/ocis/storage/ocm/ocminvites.json | Path to the JSON file where OCM invite data will be stored. This file is maintained by the instance and must not be changed manually. If not defined, the root directory derives from $OCIS_BASE_DATA_PATH/storage/ocm.|
| OCM_OCM_INVITE_MANAGER_TOKEN_EXPIRATION | Duration | 24h0m0s | Expiry duration for invite tokens.|
| OCM_OCM_INVITE_MANAGER_TIMEOUT | Duration | 30s | Timeout specifies a time limit for requests made to OCM endpoints.|
| OCM_OCM_INVITE_MANAGER_INSECURE | bool | false | Disable TLS certificate validation for the OCM connections. Do not set this in production environments.|
| SHARING_OCM_PROVIDER_AUTHORIZER_DRIVER | string | json | Driver to be used to persist ocm invites. Supported value is only 'json'.|
| OCM_OCM_PROVIDER_AUTHORIZER_PROVIDERS_FILE | string | /etc/ocis/ocmproviders.json | Path to the JSON file where ocm invite data will be stored. Defaults to $OCIS_CONFIG_DIR/ocmproviders.json.|
| OCM_OCM_SHARE_PROVIDER_DRIVER | string | json | Driver to be used for the OCM share provider. Supported value is only 'json'.|
| OCM_OCM_SHAREPROVIDER_JSON_FILE | string | /var/lib/ocis/storage/ocm/ocmshares.json | Path to the JSON file where OCM share data will be stored. If not defined, the root directory derives from $OCIS_BASE_DATA_PATH/storage.|
| OCM_OCM_SHARE_PROVIDER_INSECURE | bool | false | Disable TLS certificate validation for the OCM connections. Do not set this in production environments.|
| OCM_WEBAPP_TEMPLATE | string |  | Template for the webapp url.|
| OCM_OCM_CORE_DRIVER | string | json | Driver to be used for the OCM core. Supported value is only 'json'.|
| OCM_OCM_CORE_JSON_FILE | string | /var/lib/ocis/storage/ocm/ocmshares.json | Path to the JSON file where OCM share data will be stored. If not defined, the root directory derives from $OCIS_BASE_DATA_PATH/storage.|
| OCM_OCM_STORAGE_PROVIDER_INSECURE | bool | false | Disable TLS certificate validation for the OCM connections. Do not set this in production environments.|
| OCM_OCM_STORAGE_PROVIDER_STORAGE_ROOT | string | /var/lib/ocis/storage/ocm | Directory where the ocm storage provider persists its data like tus upload info files.|
| OCM_OCM_STORAGE_DATA_SERVER_URL | string | http://localhost:9280/data | URL of the data server, needs to be reachable by the data gateway provided by the frontend service or the user if directly exposed.|