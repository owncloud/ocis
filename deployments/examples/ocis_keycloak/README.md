---
document this deployment example in: docs/ocis/deployment/ocis_keycloak.md
---

Please refer to [our documentation](https://owncloud.dev/ocis/deployment/ocis_keycloak/)
for instructions on how to deploy this scenario.


## Vault

Adds vault sidecar services (`graph-vault`, `storage-users-vault`) to `ocis_keycloak`.

### Running

Uncomment in `.env` or run directly:

```bash
docker compose -f docker-compose.yml -f vault.yml up -d
```
