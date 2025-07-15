# Invitations

The invitations service provides an [Invitation Manager](https://learn.microsoft.com/en-us/graph/api/invitation-post?view=graph-rest-1.0&tabs=http) that can be used to invite external users, aka guests, to an organization.

*   Users invited via the Invitation Manager (via the libre graph API) will have the `userType="Guest"`.
*   Users belonging to the organization have the `userType="Member"`.

The corresponding CS3 API [user types](https://cs3org.github.io/cs3apis/#cs3.identity.user.v1beta1.UserType) used to reperesent this are: `USER_TYPE_GUEST` and `USER_TYPE_PRIMARY`.

## Provisioning Backends

When Infinite Scale is used via the IDM service for the user management, users are created using the `/graph/v1.0/users` endpoint via the libre graph API. For larger deployments, the Keycloak admin API can be used to provision users. In a future step, the endpoint, credentials and body might be made configurable using templates.

### Keycloak

The default and currently only available backend used to handle invitations is [Keycloak](https://www.keycloak.org/). Keycloak is an open source identity and access management (IAM) system which is also integrated by other Infinite Scale services as an authentication and authorization backend.

#### Keycloak Realm Configuration

<!--- Note that the link below must be an absolute URL and not a relative file path --->

See the [example configuration json file](https://github.com/owncloud/ocis/blob/master/services/invitations/md-sources/example-realm.json) of a Keycloak realm the backend will work with. This file includes the `invitations` client, which is relevant for this service.

To use the example json, set the `INVITATIONS_KEYCLOAK_CLIENT_ID` setting to `invitations`, though any other client ID can be configured. 

Importing this example into Keycloak will give you a realm that federates with an LDAP server, has the right
clients configured and all mappers correctly set. Be sure to set all the credentials after the import,
as they will be disabled.

The most relevant bits are the mappers for the `OWNCLOUD_ID` and `OWNCLOUD_USER_TYPE` user properties.

## Backend Configuration

After Keycloak has been configured, the invitation service needs to be configured with the following environment variables:

* `INVITATIONS_KEYCLOAK_BASE_PATH`: The URL to access Keycloak.
* `INVITATIONS_KEYCLOAK_CLIENT_ID`: The client ID of the client to use. In the above example, `invitations` is used.
* `INVITATIONS_KEYCLOAK_CLIENT_SECRET`: The client secret used to authenticate. This can be found in the Keycloak UI.
* `INVITATIONS_KEYCLOAK_CLIENT_REALM`: The realm where the client was added. In the example above, `ocis` is used.
* `INVITATIONS_KEYCLOAK_USER_REALM`: The realm where to add the users. In the example above, `ocis` is used.
* `INVITATIONS_KEYCLOAK_INSECURE_SKIP_VERIFY`: If set to true, the verification of the Keycloak HTTPS certificate is skipped. This is not recommended in production environments.

## Bridging Provisioning Delay

Consider that when a guest account has to be provisioned in an external user management, there might be a delay between creating the user and the user being available in the local Infinite Scale system.
