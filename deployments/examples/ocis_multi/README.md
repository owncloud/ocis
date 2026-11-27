# Infinite Scale Multi-Instance Deployment Example

This docker-compose file spins up 2 ocis instances connected to the same keycloak instance.
- `ocis.owncloud.test`
- `ocis.ocm.owncloud.test`

Demo User have different roles on different instances

| User | ocis.owncloud.test | ocis.ocm.owncloud.test |
| --- | --- | --- |
| admin | admin | admin |
| einstein | user | |
| katherine | space-admin | |
| marie | | user |
| moss | | admin |
| richard | | |


Users can be invited to instances they are not member of by using their exact email address in space membership or share dialog.

## Guest Invitations

This example does not start the invitations service, so inviting guests by
email (`/graph/v1.0/invitations`) returns `502`. To enable it:

1. Add an `invitations` client to the realm: a confidential client
   (`publicClient: false`) with only `serviceAccountsEnabled: true`, and give its
   service account the `realm-management` client role `manage-users`.
2. Start the service and point it at Keycloak in the `ocis` container:

   ```yaml
   OCIS_ADD_RUN_SERVICES: "invitations"
   INVITATIONS_KEYCLOAK_BASE_PATH: https://keycloak.owncloud.test
   INVITATIONS_KEYCLOAK_CLIENT_ID: invitations
   INVITATIONS_KEYCLOAK_CLIENT_SECRET: <client secret>
   INVITATIONS_KEYCLOAK_CLIENT_REALM: oCIS
   INVITATIONS_KEYCLOAK_USER_REALM: oCIS
   ```

   Do this for each instance that should be able to invite guests, using its
   own client or a shared one. Keep any services already listed in
   `OCIS_ADD_RUN_SERVICES` and append `invitations` to them (e.g.
   `"auth-app,invitations"`).
3. Give the realm a working SMTP server (Realm settings → Email). Keycloak sends
   the invitation email. If it can't, the invitation fails and the guest is
   removed again.

The `ocis_full` example has all of this set up when Keycloak is enabled. See the
[invitations service README](../../../services/invitations/README.md) for all
settings.

## Master-ID Configuration

The master-id grants users access to all instances without maintaining instance-specific IDs. Users with the master-id in their `owncloudMemberOf` or `owncloudGuestOf` LDAP attributes can login to any instance.

**Configuration:**
- Set `OCIS_MULTI_INSTANCE_MASTER_ID` environment variable on each instance
- Configure users in LDAP with the master-ID value in their `owncloudMemberOf` attribute
- The master-ID is automatically injected into LDAP user queries

**Example:**
```yaml
OCIS_MULTI_INSTANCE_MASTER_ID: "11111111-1111-1111-1111-111111111111"
OCIS_LDAP_USER_FILTER: "(&(objectclass=owncloud)(ownCloudMemberOf=instance-id))"
```

In this example, the admin user (configured in `config/ldap/ldif/20_users.ldif` with `owncloudMemberOf: 11111111-1111-1111-1111-111111111111`) can access both instances. The master-ID is automatically included in LDAP queries alongside the per-instance filter.

**Security Note:** Master-ID grants member status only (not admin privileges). Regular users require specific instance IDs for access control.



## Vault Mode

You can enable the vault storage on both instances. The file `vault-storage.yml` adds one extra `storage-users` service to each instance: `storage-users-vault` for `ocis` and `storage-users-vault-ocm` for `ocis-ocm`. It also enables MFA and vault mode on both `ocis` and `ocis-ocm`.

To start the deployment with the vault, add the file `vault-storage.yml` to the compose command:

```bash
docker compose -f docker-compose.yml -f vault-storage.yml up -d
```

Alternatively, set the `COMPOSE_FILE` environment variable:

```bash
COMPOSE_FILE=docker-compose.yml:vault-storage.yml docker compose up -d
```
