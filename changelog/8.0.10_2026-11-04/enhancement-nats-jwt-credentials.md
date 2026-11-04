Enhancement: Add NATS JWT authentication with a .creds file

A NATS server in operator (JWT) mode authenticates a client with a user JWT and
an NKey seed. Before this change, oCIS could connect to NATS only with a
username and a password. Thus oCIS could not use a shared NATS cluster with one
NATS account for each instance.

oCIS now accepts JWT credentials for all NATS connections: the events, the
caches, the persistent stores, the service registry, and the health checks.
You can give the path of a NATS `.creds` file or the content of such a file:

- `OCIS_NATS_CREDS_FILE`: the path of a `.creds` file. oCIS reads the file
  again on each reconnect, so a rotated JWT works without a restart.
- `OCIS_NATS_CREDS`: the content of a `.creds` file. Use it when you cannot
  mount a file.

Each service also has its own variables, for example `PROXY_EVENTS_CREDS_FILE`
and `GRAPH_CACHE_CREDS`. The service registry reads
`MICRO_REGISTRY_AUTH_CREDS_FILE` and `MICRO_REGISTRY_AUTH_CREDS`. If they are
empty, the registry uses the `OCIS_NATS_*` variables. The registry reads only
environment variables, not the `ocis.yaml` file.

A service-specific variable takes precedence over a global variable. In the
same scope, the content takes precedence over the file.

When JWT credentials are set, oCIS does not send the username and the
password. This rule also uses the scope. Example: `OCIS_NATS_CREDS_FILE` is set,
and `GRAPH_CACHE_AUTH_USERNAME` and `GRAPH_CACHE_AUTH_PASSWORD` are set. Then
the graph cache uses the username and the password. A username or a password
alone does not replace the JWT credentials. The same applies to `MICRO_REGISTRY_AUTH_USERNAME` and
`MICRO_REGISTRY_AUTH_PASSWORD` for the service registry.

oCIS checks the credentials at startup and gives a clear error for an invalid
value. oCIS checks the credentials of a store or a cache only when its type is
`nats-js` or `nats-js-kv`. Other types do not use the credentials.

The embedded nats service does not support JWT mode. To use JWT
authentication, run an external NATS cluster in operator mode and do not start
the embedded nats service, for example with `OCIS_EXCLUDE_RUN_SERVICES=nats`.

The NATS health checks of the services now use the same authentication and TLS
settings as the service connections. Before this change, the health checks
connected without credentials and without TLS. The proxy health check also used
the events cluster ID as the NATS address. It now uses the events endpoint.
Each health check skips the check when the NATS address is empty. Each health
check keeps one NATS connection and uses it again for each probe. Each probe
sends a PING to the NATS server, so the check finds a lost server at once. A
slow NATS server does not block the other probes. The health check of a store
succeeds when one node of the store is reachable, like the store connection.

The userlog service had no default node for its store. With
`OCIS_PERSISTENT_STORE=nats-js-kv` it could not connect to NATS. It now uses
`127.0.0.1:9233` like the other services.

The create home cache of the gateway service did not get the TLS settings. It
now has the variables `GATEWAY_CREATE_HOME_CACHE_ENABLE_TLS`,
`GATEWAY_CREATE_HOME_CACHE_TLS_INSECURE`, and
`GATEWAY_CREATE_HOME_CACHE_TLS_ROOT_CA_CERTIFICATE`, and it uses the global
`OCIS_CACHE_*` TLS variables.

https://github.com/owncloud/ocis/pull/13220
