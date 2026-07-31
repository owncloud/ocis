Enhancement: Provision invited guests directly into the identity backend

The invitations service can now provision invited guests directly into the oCIS
identity backend (LDAP) instead of an external Keycloak IdP. The backend is selected
with the new `INVITATIONS_BACKEND` environment variable: `keycloak` (default, the
previous behaviour) or `ldap`.

A guest created via the `ldap` backend is written into the same directory oCIS reads
for share-recipient resolution, so it is immediately resolvable through the Graph API
and usable as a share recipient right after the invitation. This closes the
provisioning-delay gap that previously existed between creating a guest in an external
IdP and being able to share with them. The `ldap` backend requires
`OCIS_LDAP_SERVER_WRITE_ENABLED=true`.

The `ocis_full` example deployment enables this by default so invited guests can be
shared with out of the box. Note: the `ldap` backend only provisions the guest; it
does not send a credential-setup email (that remains a Keycloak feature).

https://github.com/owncloud/ocis/pull/12469
