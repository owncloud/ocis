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

## Data Storage

All persistent runtime data (oCIS config/data, S3 storage, database, virus
definitions, web app extensions, etc.) is bind-mounted from a `./data/`
directory next to this compose project, rather than docker-managed named
volumes. The oCIS config/data paths can be redirected elsewhere via the
`OCIS_CONFIG_DIR`/`OCIS_DATA_DIR` variables in `.env`.

Upgrading an existing deployment that still has data in the old named
volumes? Run `./migrate-volumes.sh` once, before starting the stack, to copy
data from the old volumes into the new `./data/` directories.

## Optional Services

### Keycloak

Keycloak can be optionally enabled by uncommenting the corresponding variables in the `.env` file:
- `KEYCLOAK=:keycloak.yml`

Note that Keycloak requires the default `ocis` Identity Provider to be disabled, which is automatically handled when the `keycloak.yml` configuration is used.

### oCIS Workflows

[oCIS Workflows](https://github.com/owncloud/ocis-workflows) is an AI-powered file workflow automation extension. It can be optionally enabled by uncommenting the corresponding variables in the `.env` file:
- `WORKFLOWS=:workflows.yml`

Note that Workflows requires the `ocis-apps` volume provided by the web extensions configuration, so `EXTENSIONS=:web_extensions/extensions.yml` must also be enabled. You also need to set `WORKFLOWS_ENCRYPTION_KEY` (e.g. via `openssl rand -base64 32`) and point `WORKFLOWS_LLM_ENDPOINT` / `WORKFLOWS_LLM_MODEL` / `WORKFLOWS_LLM_API_KEY` at an OpenAI-compatible LLM backend.

To use Workflows *automations* (scheduled or event-triggered runs, rather than only manual "Run now"), the backend mints oCIS app-passwords on the user's behalf, which requires oCIS's `auth-app` service. App-password support (`PROXY_ENABLE_APP_AUTH`) and the `auth-app` service are both off by default, since they widen the oCIS instance's authentication surface, so you must explicitly:
- Set `PROXY_ENABLE_APP_AUTH=true` in the `.env` file.
- Add `auth-app` to `START_ADDITIONAL_SERVICES` in the `.env` file, e.g. `START_ADDITIONAL_SERVICES="notifications,auth-app"`.
