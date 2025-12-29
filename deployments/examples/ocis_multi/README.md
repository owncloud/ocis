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


