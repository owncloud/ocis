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

## Optional Features

`ocis_full` enables optional components through `.env` and `COMPOSE_FILE`.

Enable Keycloak:

```bash
KEYCLOAK=:keycloak.yml
```

Enable OpenLDAP:

```bash
OPENLDAP=:openldap.yml
```

Enable both:

```bash
KEYCLOAK=:keycloak.yml
OPENLDAP=:openldap.yml
```

When both are enabled, the recommended full setup is:

- Keycloak is the OIDC provider for user authentication.
- OpenLDAP is the identity directory backend for oCIS Graph users and groups.
- Optional recommendation: configure Keycloak User Federation to read users from LDAP.

## Run Commands

```bash
docker compose pull
docker compose up -d
docker compose ps
docker compose logs -f ocis
docker compose logs -f keycloak
docker compose logs -f openldap
```

## Troubleshooting

- OIDC issuer / redirect URI mismatch:
  Confirm `OCIS_OIDC_ISSUER` and `KEYCLOAK_DOMAIN` match the Keycloak realm URL and that the Keycloak `web` client redirect URI matches your `OCIS_DOMAIN`.
- Keycloak behind reverse proxy:
  Verify `KC_PROXY_HEADERS=xforwarded`, `KC_HTTP_ENABLED=true`, and Traefik is routing `Host(<KEYCLOAK_DOMAIN>)` to Keycloak port `8080`.
- LDAP bind DN / base DN issues:
  Check `LDAP_ADMIN_PASSWORD`, bind DN `cn=admin,dc=owncloud,dc=com`, user base DN `ou=users,dc=owncloud,dc=com`, and group base DN `ou=groups,dc=owncloud,dc=com`.
- oCIS cannot list users/groups:
  Inspect logs from `ocis` and `openldap`, then verify `GRAPH_LDAP_*` settings and whether `GRAPH_LDAP_SERVER_WRITE_ENABLED` should remain `false` (read-only) or be switched to `true`.
