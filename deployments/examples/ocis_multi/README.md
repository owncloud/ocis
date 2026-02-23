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

**Configuration:**
- **Master-ID Value**: `11111111-1111-1111-1111-111111111111`
- **Environment Variable**: `OCIS_MULTI_INSTANCE_MASTER_ID`
- **User**: admin (configured in config/ldap/ldif/20_users.ldif)
- **Purpose**: Users with this ID in their `owncloudMemberOf` or `owncloudGuestOf` claims can login to any instance

Users with the master-id are granted member access to all instances without maintaining instance-specific IDs. The feature is optional and disabled when `OCIS_MULTI_INSTANCE_MASTER_ID` is empty.

**Security Note:** Master-id grants member status only (not admin privileges). Regular users require specific instance IDs for access control.


