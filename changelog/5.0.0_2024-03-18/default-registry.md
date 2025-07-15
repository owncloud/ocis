Enhancement: Make nats-js-kv the default registry

The previously used default `mdns` is faulty. Deprecated it together with `consul`, `nats` and `etcd` implementations.

https://github.com/owncloud/ocis/pull/8011
https://github.com/owncloud/ocis/pull/8027
