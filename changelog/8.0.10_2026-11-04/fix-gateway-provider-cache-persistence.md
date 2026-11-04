Bugfix: Disable the persistence of the gateway provider cache

The gateway service sent the wrong config key to reva for the provider cache.
Reva reads the key `cache_disable_persistence`, but the gateway sent
`disable_persistence`. Thus the variables `OCIS_CACHE_DISABLE_PERSISTENCE` and
`GATEWAY_PROVIDER_CACHE_DISABLE_PERSISTENCE` had no effect on the provider
cache. The provider cache always used persistence.

The gateway now sends the correct key. Check your configuration before you
upgrade. If `OCIS_CACHE_DISABLE_PERSISTENCE=true` is set, the provider cache
now stops the persistence. This is the documented behavior.

https://github.com/owncloud/ocis/pull/13220
