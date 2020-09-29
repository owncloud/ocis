Enhancement: Add basic auth option

We added a new `enable-basic-auth` option and `PROXY_ENABLE_BASIC_AUTH` environment variable that can be set to `true` to make the proxy verify the basic auth header with the accounts service. This should only be used for testing and development and is disabled by default.

https://github.com/owncloud/ocis/pull/627
https://github.com/owncloud/product/issues/198
