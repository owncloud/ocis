Change: use glauth as ldap backend, default to running behind ocis-proxy

We changed the default configuration to integrate better with ocis.

The default ldap port changes to 9125, which is used by ocis-glauth and we use ocis-proxy to do the tls offloading.
Clients are supposed to use the ocis-proxy endpoint `https://localhost:9200`

https://github.com/owncloud/ocis-konnectd/pull/52
