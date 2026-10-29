# Nats

The nats service is the event broker of the system. It distributes events among all other services and enables other services to communicate asynchronous.

Services can `Publish` events to the nats service and nats will store these events on disk and distribute these events to other services eventually. Services can `Consume` events from the nats service by registering to a `ConsumerGroup`. Each `ConsumerGroup` is guaranteed to get each event exactly once. In most cases, each service will register its own `ConsumerGroup`. When there are multiple instances of a service, those instances will usually use that `ConsumerGroup` as common resource.

## Underlying Technology

As the service name suggests, this service is based on [NATS](https://nats.io/) specifically on [NATS Jetstream](https://docs.nats.io/nats-concepts/jetstream) to enable persistence.

## Default Registry

By default, `nats-js-kv` is configured as embedded default registry via the `MICRO_REGISTRY` environment variable. If you do not want using the build-in nats registry, set `MICRO_REGISTRY_ADDRESS` to the address of the nats-js cluster, which is the same value as `OCIS_EVENTS_ENDPOINT`. Optionally use `MICRO_REGISTRY_AUTH_USERNAME` and `MICRO_REGISTRY_AUTH_PASSWORD` to authenticate with the external nats cluster.

## Persistance

To be able to deliver events even after a system or service restart, nats will store events in a folder on the local filesystem. This folder can be specified by setting the `NATS_NATS_STORE_DIR` enviroment variable. If not set, the service will fall back to `$OCIS_BASE_DATA_PATH/nats`.

## TLS Encryption

Connections to the nats service (`Publisher`/`Consumer` see above) can be TLS encrypted by setting the corresponding env vars `NATS_TLS_CERT`, `NATS_TLS_KEY` to the cert and key files and `ENABLE_TLS` to true. Checking the certificate of incoming request can be disabled with the `NATS_EVENTS_ENABLE_TLS` environment variable.

Certificate files can also be set via global variables starting with `OCIS_`, for details see the environment variable list.

Note that using TLS is highly recommended for productive environments, especially when using container orchestration with Kubernetes.

## JWT Authentication With an External NATS Cluster

The services can authenticate with an external NATS cluster that runs in operator (JWT) mode. Use this when several oCIS instances share one NATS cluster. Each instance then uses its own NATS account. An account has its own subjects, its own JetStream data, and its own limits.

In JWT mode, a client authenticates with a user JWT and an NKey seed. A NATS `.creds` file holds both. Set one of these environment variables:

* `OCIS_NATS_CREDS_FILE`: the path to the `.creds` file. The services read the file again on each reconnect, so a rotated JWT works without a restart.
* `OCIS_NATS_CREDS`: the content of the `.creds` file. Use this variable when you cannot mount a file. A rotated JWT then needs a restart.

When you set more than one of these variables, these rules apply:

* A service-specific variable takes precedence over a global variable. Example: you set `OCIS_NATS_CREDS` and `GRAPH_CACHE_CREDS_FILE`. Then the graph cache uses the file of `GRAPH_CACHE_CREDS_FILE`.
* In the same scope, the content takes precedence over the file. Example: you set `OCIS_NATS_CREDS` and `OCIS_NATS_CREDS_FILE`. Then the services use the content of `OCIS_NATS_CREDS`.

The services use these credentials for the events, the caches, the persistent stores, and the health checks. When the credentials are set, the services do not send a username and a password. Exception: a service-specific username or password takes precedence over the global variables. Example: you set `OCIS_NATS_CREDS_FILE` and `GRAPH_CACHE_AUTH_USERNAME`. Then the graph cache uses the username and the password. Each service also has its own variables, for example `PROXY_EVENTS_CREDS_FILE`. For details, see the environment variable list.

The service registry reads only environment variables. It reads `MICRO_REGISTRY_AUTH_CREDS_FILE` and `MICRO_REGISTRY_AUTH_CREDS`. If these variables are empty, it uses `OCIS_NATS_CREDS_FILE` and `OCIS_NATS_CREDS`. The same rules apply: the registry variables take precedence over the global variables. This also applies to `MICRO_REGISTRY_AUTH_USERNAME` and `MICRO_REGISTRY_AUTH_PASSWORD`. A value in the `ocis.yaml` file does not reach the registry.

The embedded nats service does not support JWT mode. It has no settings for an operator, accounts, or a resolver. In the default single binary setup, the embedded nats service and all its clients run in the same process, so JWT mode gives no benefit there. To use JWT authentication, do these steps:

1. Run an external NATS cluster in operator mode. Create the operator, the accounts, and the users with the [nsc](https://docs.nats.io/using-nats/nats-tools/nsc) tool.
2. Do not start the embedded nats service, for example with `OCIS_EXCLUDE_RUN_SERVICES=nats`.
3. Set the address of the external cluster in `OCIS_EVENTS_ENDPOINT`, `OCIS_CACHE_STORE_NODES`, `OCIS_PERSISTENT_STORE_NODES`, and `MICRO_REGISTRY_ADDRESS`.
4. Set `OCIS_NATS_CREDS_FILE` or `OCIS_NATS_CREDS`.
