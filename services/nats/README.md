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
