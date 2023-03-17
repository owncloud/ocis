# ocis

The ocis package contains the Infinite Scale runtime and the commands for the Infinite Scale cli.

## Service registry

This package also configures the service registry which will be used to look up the service addresses. It defaults to mDNS. Keep that in mind when using systems with mDNS disabled by default (i.e. SUSE).

Available registries are:

-   nats
-   kubernetes
-   etcd
-   consul
-   memory
-   mdns (default)

To configure which registry to use, you have to set the environment variable `MICRO_REGISTRY`, and for all except `memory` and `mdns` you also have to set the registry address via `MICRO_REGISTRY_ADDRESS`.

### etcd

To authenticate the connection to the etcd registry, you have to set `ETCD_USERNAME` and `ETCD_PASSWORD`.

## Memory limits

oCIS will automatically set the go native `GOMEMLIMIT` to `0.9`. To disable the limit set `AUTOMEMEMLIMIT=off`. For more information take a look at the official [Guide to the Go Garbage Collector](https://go.dev/doc/gc-guide).