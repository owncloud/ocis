Bugfix: Remove the Keycloak guest when sending the invitation mail fails

When the invitations service created a guest in Keycloak but sending the
execute-actions email failed, for example because a configured required action
is not enabled in the realm, the guest stayed behind in Keycloak and every
retry failed with a conflict. The invitations service now tries to delete the
just-created guest if Keycloak refuses to send the mail and returns the original
error. If the outcome is unknown, for example after a timeout, the guest is kept
because the mail may have been sent. If the deletion fails, it is logged with
the email address and user id and the guest has to be removed manually, while
the original error is still returned. The Keycloak service account of the
invitations client needs the `manage-users` role for this, which is now
documented together with the `INVITATIONS_KEYCLOAK_EXECUTE_ACTIONS` setting.

https://github.com/owncloud/ocis/issues/13253
https://github.com/owncloud/ocis/pull/13408
