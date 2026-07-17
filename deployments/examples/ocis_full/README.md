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

## Optional Services

### Keycloak

Keycloak can be optionally enabled by uncommenting the corresponding variables in the `.env` file:
- `KEYCLOAK=:keycloak.yml`

Note that Keycloak requires the default `ocis` Identity Provider to be disabled, which is automatically handled when the `keycloak.yml` configuration is used.

### MCP Server

The [oCIS MCP Server](https://github.com/owncloud/ocis-mcp-server) exposes this deployment as a set of AI-accessible tools over the Model Context Protocol (MCP), so AI assistants such as Claude can manage users, spaces, files and shares through natural language. It can be optionally enabled by uncommenting the corresponding variable in the `.env` file:
- `OCIS_MCP=:mcp.yml`

It is reachable at `https://${OCIS_MCP_DOMAIN:-mcp.owncloud.test}` once enabled. Two things must be configured in `.env` before it will start:
- `OCIS_MCP_HTTP_SECRET` — a shared secret every MCP client must send as `Authorization: Bearer <secret>`. The container refuses to start without it.
- `OCIS_MCP_APP_TOKEN_USER` / `OCIS_MCP_APP_TOKEN_VALUE` — an oCIS app token, created with `docker-compose exec ocis ocis auth-app create --user-name="admin" --expiration="8760h"`.
