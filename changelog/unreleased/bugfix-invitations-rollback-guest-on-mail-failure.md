Bugfix: Remove the Keycloak guest when sending the invitation mail fails

When the invitations service created a guest in Keycloak but sending the
execute-actions email failed, for example because a configured required action
is not enabled in the realm, the guest stayed behind in Keycloak and every
retry failed with a conflict. The invitations service now deletes the
just-created guest if the mail can not be sent and returns the original error.
If the deletion fails as well, it is logged and the original error is still
returned. The Keycloak service account of the invitations client needs the
`manage-users` role for this, which is now documented together with the
`INVITATIONS_KEYCLOAK_EXECUTE_ACTIONS` setting.

https://github.com/owncloud/ocis/issues/13253
