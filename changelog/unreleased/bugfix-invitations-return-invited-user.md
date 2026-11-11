Bugfix: Return the created user in the invitation response

The invitations service created the guest account but returned the original request
body unchanged, so the `invitedUser` relation defined by the Graph invitation
response was always empty. Clients had no way to learn the identity of the user that
was just created. The service now populates `invitedUser` with the created user (id,
email address and user type), where the id is the OWNCLOUD_ID that oCIS uses once the
guest is provisioned locally.

https://github.com/owncloud/ocis/pull/12467
