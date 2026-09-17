Enhancement: Configurable Keycloak invitation actions and working ocis_full SMTP

The invitations service now lets operators configure the Keycloak required actions
that are sent to invited guests via the execute-actions email, using the new
`INVITATIONS_KEYCLOAK_EXECUTE_ACTIONS` environment variable. It defaults to
`UPDATE_PASSWORD,VERIFY_EMAIL`; configured values are passed to Keycloak as-is, and
the service falls back to the defaults when none are configured, so an invited guest
always has a way to set up their account.

The Keycloak realm shipped with the `ocis_full` example deployment now points its
SMTP server at the bundled mailpit instance. Previously the realm shipped with an
empty SMTP configuration, so Keycloak silently dropped every mail it tried to send
(including the guest invitation email) even though the invitations service
correctly requested it.

https://github.com/owncloud/ocis/pull/12467
