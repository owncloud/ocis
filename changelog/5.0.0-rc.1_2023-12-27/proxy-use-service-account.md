Enhancement: Proxy uses service accounts for provisioning

The proxy service now uses a service account for provsioning task, like role
assignment and user auto-provisioning. This cleans up some technical debt that
required us to mint reva tokes inside the proxy service.

https://github.com/owncloud/ocis/pull/7240
https://github.com/owncloud/ocis/issues/5550
