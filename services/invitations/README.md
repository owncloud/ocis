# Invitations Service

The invitations service provides an [Invitation Manager](https://learn.microsoft.com/en-us/graph/api/invitation-post?view=graph-rest-1.0&tabs=http) that can be used to invite external users aka Guests to an organization.

Users invited via this Invitation Manager (libre graph API) will have `userType="Guest"`, whereas users belonging to the organization have `userType="Member"`.

The corresponding CS3 API [user types](https://cs3org.github.io/cs3apis/#cs3.identity.user.v1beta1.UserType) used to reperesent this are: `USER_TYPE_GUEST` and `USER_TYPE_PRIMARY`.


## Provisioning backends

When oCIS is used for user management the users are created using the `/graph/v1.0/users` endpoint. For larger deployments the keycloak admin API can be used to provision users. We might even make the endpoint, credentials and body configurable using templates.


## Bridging provisioning delay

When a guest account has to be provisioned in an external user management there might be a delay between creating the user and it being available in the local ocis system. In the first iteration the invitations service will only keep track of invites in memory. This list could be persisted in future iterations.
