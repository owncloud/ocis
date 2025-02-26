Enhancement: Add initial nats and kubernetes registry support

We added initial support to use nats and kubernetes as a service registry using `MICRO_REGISTRY=nats` and `MICRO_REGISTRY=kubernetes` respectively.
Multiple nodes can be given with `MICRO_REGISTRY_ADDRESS=1.2.3.4,5.6.7.8,9.10.11.12`.

https://github.com/owncloud/ocis/pull/1697
