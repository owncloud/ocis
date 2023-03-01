# Invitations Service

The invitations service provides an [Invitation Manager](https://learn.microsoft.com/en-us/graph/api/invitation-post?view=graph-rest-1.0&tabs=http) that can be used to invide external users aka Guests to an organization.

On the libre graph API invited users will have `userType="Guest"`, whereas users belonging to the organization have `userType="Member"`.

The corresponding CS3 API [user types](https://cs3org.github.io/cs3apis/#cs3.identity.user.v1beta1.UserType) used to reperesent this are USER_TYPE_GUEST and USER_TYPE_PRIMARY.