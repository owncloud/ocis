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

### Kiteworks Storage

Enables a read-only Kiteworks storage provider alongside the default oCIS storage.
Kiteworks top-level folders appear as project spaces in the oCIS web UI.
Creating new project spaces from the UI continues to use decomposedfs.

**Prerequisites:**

The Kiteworks storage driver is not part of an upstream oCIS release. Build the `owncloud/ocis:dev` image locally from this branch:
```bash
make -C ocis dev-docker
```

**Required `.env` variables:**
```env
KITEWORKS_ENDPOINT=https://your-kiteworks-instance.example.com
KITEWORKS_API_TOKEN=<your-api-token>
```

**Optional `.env` variables:**
```env
KITEWORKS_MOUNT_ID=aa309def-b364-417b-8e41-85b8a9393c4b  # stable UUID; change only to avoid collision with another deployment
KITEWORKS_INSECURE=false                                   # set true only if the Kiteworks TLS cert is self-signed
```

**Start command:**
```bash
docker compose -f docker-compose.yml -f ocis.yml -f kiteworks-storage.yml up -d
```

**Verify:**
```bash
# Kiteworks folders appear as project spaces alongside the personal decomposedfs space
curl -su admin:admin "https://${OCIS_DOMAIN}/graph/v1.0/me/drives" \
  | jq '.value[] | {name, driveType}'

# Creating a new project space must succeed (HTTP 201) and land on decomposedfs
curl -su admin:admin -X POST "https://${OCIS_DOMAIN}/graph/v1.0/drives" \
  -H "Content-Type: application/json" \
  -d '{"name":"my space","driveType":"project"}' -w "\nHTTP %{http_code}\n"
```
