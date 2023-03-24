# Invitations Service

The invitations service provides an [Invitation Manager](https://learn.microsoft.com/en-us/graph/api/invitation-post?view=graph-rest-1.0&tabs=http) that can be used to invite external users aka Guests to an organization.

Users invited via this Invitation Manager (libre graph API) will have `userType="Guest"`, whereas users belonging to the organization have `userType="Member"`.

The corresponding CS3 API [user types](https://cs3org.github.io/cs3apis/#cs3.identity.user.v1beta1.UserType) used to reperesent this are: `USER_TYPE_GUEST` and `USER_TYPE_PRIMARY`.


## Commands

Right now, the invitations server supplies the following commands:

```
COMMANDS:
   server   start the invitations service without runtime (unsupervised mode)
   health   check health status
   version  print the version of this binary and the running service instances
   help, h  Shows a list of commands or help for one command
```

### Command examples

**Start the invitiations service**

```sh
$ ocis invitations server
```

## Provisioning Backends

### Keycloak

The default and currently only available backend used to handle Invitations is [Keycloak](https://www.keycloak.org/). Keycloak is an open source identity and access management (IAM) system which is also integrated by other OCIS services as an authentication and authorization backend.

#### Realm Configuration

We supply an [example configuration json file](https://github.com/owncloud/ocis/blob/master/services/invitations/examples/keycloak/example-realm.json) of a Keycloak realm that the backend will work with. This file includes the `invitations` client, which is relevant for this service.

Importing this into keycloak will give you a realm that federates with an LDAP server, has the right
clients configured, and all the mappers correctly set. Be sure to set all the credentials after the import,
as they will be disabled.

The most relevant bits here are the mappers for the `OWNCLOUD_ID` and `OWNCLOUD_USER_TYPE` user properties.

#### Backend Configuration

After Keycloak has been configured, the invitation service needs to be configured with the following environment variables:

* `INVITATIONS_KEYCLOAK_BASE_PATH`: The URL to access Keycloak.
* `INVITATIONS_KEYCLOAK_CLIENT_ID`: The client ID of the client to use. In the above example, `invitations` is used.
* `INVITATIONS_KEYCLOAK_CLIENT_SECRET`: The client secret used to authenticate. This can be found in the Keycloak UI.
* `INVITATIONS_KEYCLOAK_CLIENT_REALM`: The realm where the client was added. In the example above, `ocis` is used.
* `INVITATIONS_KEYCLOAK_USER_REALM`: The realm where to add the users. In the example above, `ocis` is used.
* `INVITATIONS_KEYCLOAK_INSECURE_SKIP_VERIFY`: If set to true, the verification of the Keycloak https certificate is skipped. This is not recommended in production enviroments.
