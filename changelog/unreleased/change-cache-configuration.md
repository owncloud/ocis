Change: Updated Cache Configuration

We updated all cache related environment vars to more closely follow the go micro naming pattern:
- `{service}_CACHE_STORE_TYPE` becomes `{service}_CACHE_STORE` or `{service}_PERSISTENT_STORE`
- `{service}_CACHE_STORE_ADDRESS(ES)` becomes `{service}_CACHE_STORE_NODES`
- The `mem` store implementation name changes to `memory`
- In yaml files the cache `type` becomes `store`
We introduced `redis-sentinel` as a store implementation.

https://github.com/owncloud/ocis/pull/5829
