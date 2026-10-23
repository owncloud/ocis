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
  mount a file. This value takes precedence over `OCIS_NATS_CREDS_FILE`.

Each service also has its own variables, for example `PROXY_EVENTS_CREDS_FILE`
and `GRAPH_CACHE_CREDS`. The service registry reads
`MICRO_REGISTRY_AUTH_CREDS_FILE` and `MICRO_REGISTRY_AUTH_CREDS`. If they are
empty, the registry uses the `OCIS_NATS_*` variables. The registry reads only
environment variables, not the `ocis.yaml` file.

When JWT credentials are set, oCIS does not send the username and the
password. oCIS checks the credentials at startup and gives a clear error for an
invalid value.

The NATS health checks of the services now use the same authentication and TLS
settings as the service connections. Before this change, the health checks
connected without credentials and without TLS.

https://github.com/owncloud/ocis/pull/13219
