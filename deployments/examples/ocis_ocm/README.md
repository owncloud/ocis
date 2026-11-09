# Infinite Scale OCM Deployment Example

This docker-compose file spins up 3 ocis instances connected via ocm.
- `ocis.owncloud.test` has a onlyoffice instance as weboffice
- `ocis.ocm.owncloud.test` has a collabora instance as weboffice
- `ocis.ocm.next.owncloud.test` has no weboffice

All instances create demo users by default.

## Collabora

The Collabora instance uses the following recommended configuration:

- `COLLABORA_ADMIN_USER` and `COLLABORA_ADMIN_PASSWORD` must be set. There are no defaults.
- Collabora runs on its own internal network (`collabora-ocm-net`) that is only shared with the proxy.
- `config/coolwsd/coolwsd.xml` is mounted into the Collabora container. It restricts the hosts Collabora may
  talk to and pins the accepted WOPI host to `wopiserver.ocm.owncloud.test`. If you change
  `WOPISERVER_OCM_DOMAIN`, adjust the hosts in this file accordingly.
