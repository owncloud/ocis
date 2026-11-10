---
document this deployment example in: docs/ocis/deployment/ocis_full.md
---

# Infinite Scale WOPI Deployment Example

This deployment example is documented in two locations for different audiences:

* In the [Admin Documentation](https://doc.owncloud.com/ocis/latest/index.html)\
  Providing two variants using detailed configuration step by step guides:\
  [Local Production Setup](https://doc.owncloud.com/ocis/next/depl-examples/ubuntu-compose/ubuntu-compose-prod.html) and [Deploy Infinite Scale on the Hetzner Cloud](https://doc.owncloud.com/ocis/next/depl-examples/ubuntu-compose/ubuntu-compose-hetzner.html).\
  Note that these examples use LetsEncrypt certificates and are intended for production use.

* In the [Developer Documentation](https://owncloud.dev/ocis/deployment/ocis_full/)\
  Providing details which are more developer focused. This description can also be used when deviating from the default.\
  Note that this examples uses self signed certificates and is intended for testing purposes.

## Collabora

When Collabora is enabled (`COLLABORA=:collabora.yml`), the example uses the following recommended configuration:

- `COLLABORA_ADMIN_USER` and `COLLABORA_ADMIN_PASSWORD` must be set in `.env`. There are no defaults.
- Collabora runs on its own internal network (`collabora-net`) that is only shared with the proxy and the
  collaboration service. It has no access to the other services of the example.
- `config/coolwsd/coolwsd.xml` is mounted into the Collabora container. It restricts the hosts Collabora may
  talk to (the WOPI host, i.e. the collaboration service) and pins the accepted WOPI host. If you change the
  WOPI source of the collaboration service, adjust the host in this file accordingly.
- The file is based on the default configuration of the Collabora image tag set in `COLLABORA_DOCKER_TAG`.
  When you update the tag, compare it with the default `/etc/coolwsd/coolwsd.xml` of the new image and merge
  changes.

## Optional Services

### Keycloak

Keycloak can be optionally enabled by uncommenting the corresponding variables in the `.env` file:
- `KEYCLOAK=:keycloak.yml`

Note that Keycloak requires the default `ocis` Identity Provider to be disabled, which is automatically handled when the `keycloak.yml` configuration is used.
