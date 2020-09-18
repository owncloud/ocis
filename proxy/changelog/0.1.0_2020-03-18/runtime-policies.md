Enhancement: Load Proxy Policies at Runtime

While a proxy without policies is of no use, the current state of ocis-proxy expects a config file either at an expected Viper location or specified via -- config-file flag.
To ease deployments and ensure a working set of policies out of the box we need a series of defaults.

https://github.com/owncloud/ocis-proxy/issues/17
https://github.com/owncloud/ocis-proxy/pull/16
