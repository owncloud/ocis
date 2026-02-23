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

## Master-ID Configuration

The master-id grants users access to all instances without maintaining instance-specific IDs. Users with the master-id in their `owncloudMemberOf` or `owncloudGuestOf` LDAP attributes can login to any instance.

**Configuration:**
- Set `OCIS_MULTI_INSTANCE_MASTER_ID` environment variable on each instance
- Configure users in LDAP with the master-id value in their `owncloudMemberOf` attribute
- The master-id is automatically injected into LDAP user queries

**Example:**
```yaml
OCIS_MULTI_INSTANCE_MASTER_ID: "11111111-1111-1111-1111-111111111111"
OCIS_LDAP_USER_FILTER: "(&(objectclass=owncloud)(ownCloudMemberOf=instance-id))"
```

In this example, the admin user (configured in `config/ldap/ldif/20_users.ldif` with `owncloudMemberOf: 11111111-1111-1111-1111-111111111111`) can access both instances. The master-id is automatically included in LDAP queries alongside the per-instance filter.

**Security Note:** Master-id grants member status only (not admin privileges). Regular users require specific instance IDs for access control.


