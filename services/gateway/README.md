# Gateway

The gateway service is responsible for passing requests to the storage providers. Other services never talk to the storage providers directly but will always send their requests via the `gateway` service.

## Caching

The gateway service is using caching as it is highly frequented with the same requests. As of now it uses two different caches:
  -   the `provider cache` is caching requests to list or get storage providers.
  -   the `create home cache` is caching requests to create personal spaces (as they only need to be executed once).

Both caches can be configured via the `OCIS_CACHE_*` envvars (or `GATEWAY_PROVIDER_CACHE_*` and `GATEWAY_CREATE_HOME_CACHE_*` respectively). See the [envvar section](/services/gateway/configuration/#environment-variables) for details.

Use `OCIS_CACHE_STORE` (`GATEWAY_PROVIDER_CACHE_STORE`, `GATEWAY_CREATE_HOME_CACHE_STORE`) to define the type of cache to use:
  -   `memory`: Basic in-memory store and the default.
  -   `redis-sentinel`: Stores data in a configured Redis Sentinel cluster.
  -   `nats-js-kv`: Stores data using key-value-store feature of [nats jetstream](https://docs.nats.io/nats-concepts/jetstream/key-value-store)
  -   `noop`: Stores nothing. Useful for testing. Not recommended in production environments.
  -   `ocmem`: Advanced in-memory store allowing max size. (deprecated)
  -   `redis`: Stores data in a configured Redis cluster. (deprecated)
  -   `etcd`: Stores data in a configured etcd cluster. (deprecated)
  -   `nats-js`: Stores data using object-store feature of [nats jetstream](https://docs.nats.io/nats-concepts/jetstream/obj_store) (deprecated)

Other store types may work but are not supported currently.

Note: The service can only be scaled if not using `memory` store and the stores are configured identically over all instances!

Note that if you have used one of the deprecated stores, you should reconfigure to one of the supported ones as the deprecated stores will be removed in a later version.

Store specific notes:
  -   When using `redis-sentinel`, the Redis master to use is configured via e.g. `OCIS_CACHE_STORE_NODES` in the form of `<sentinel-host>:<sentinel-port>/<redis-master>` like `10.10.0.200:26379/mymaster`.
  -   When using `nats-js-kv` it is recommended to set `OCIS_CACHE_STORE_NODES` to the same value as `OCIS_EVENTS_ENDPOINT`. That way the cache uses the same nats instance as the event bus.
  -   When using the `nats-js-kv` store, it is possible to set `OCIS_CACHE_DISABLE_PERSISTENCE` to instruct nats to not persist cache data on disc.

## Storage registry

In order to add another storage provider the CS3 storage registry that is running as part of the CS3 gateway hes to be made aware of it. The easiest cleanest way to do it is to set `GATEWAY_STORAGE_REGISTRY_CONFIG_JSON=/path/to/storages.json` and list all storage providers like this:

```json
{
  "com.owncloud.api.storage-users": {
    "providerid": "{storage-users-mount-uuid}",
    "spaces": {
      "personal": {
        "mount_point":   "/users",
        "path_template": "/users/{{.Space.Owner.Id.OpaqueId}}"
      },
      "project": {
        "mount_point":   "/projects",
        "path_template": "/projects/{{.Space.Name}}"
      }
    }
  },
  "com.owncloud.api.storage-shares": {
    "providerid": "a0ca6a90-a365-4782-871e-d44447bbc668",
    "spaces": {
      "virtual": {
        "mount_point": "/users/{{.CurrentUser.Id.OpaqueId}}/Shares"
      },
      "grant": {
        "mount_point": "."
      },
      "mountpoint": {
        "mount_point":   "/users/{{.CurrentUser.Id.OpaqueId}}/Shares",
        "path_template": "/users/{{.CurrentUser.Id.OpaqueId}}/Shares/{{.Space.Name}}"
      }
    }
  },
  "com.owncloud.api.storage-publiclink": {
    "providerid": "7993447f-687f-490d-875c-ac95e89a62a4",
    "spaces": {
      "grant": {
        "mount_point": "."
      },
      "mountpoint": {
        "mount_point":   "/public",
        "path_template": "/public/{{.Space.Root.OpaqueId}}"
      }
    }
  },
  "com.owncloud.api.ocm": {
    "providerid": "89f37a33-858b-45fa-8890-a1f2b27d90e1",
    "spaces": {
      "grant": {
        "mount_point": "."
      },
      "mountpoint": {
        "mount_point":   "/ocm",
        "path_template": "/ocm/{{.Space.Root.OpaqueId}}"
      }
    }
  },
  "com.owncloud.api.storage-hello": {
    "providerid": "hello-storage-id",
    "spaces": {
      "project": {
        "mount_point":   "/hello",
        "path_template": "/hello/{{.Space.Name}}"
      }
    }
  }
}
```

In the above replace `{storage-users-mount-uuid}` with the mount UUID that was generated for the storage-users service. You can find it in the `config.yaml` generated on by `ocis init`. The last entry `com.owncloud.api.storage-hello` and its `providerid` `"hello-storage-id"` are an example for in additional storage provider, in this case running `hellofs`, an example minimal storage driver.