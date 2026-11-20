Change: MFA expires after the MFA session duration

Users now have to enter their second factor again once it is older than
`OCIS_MFA_SESSION_DURATION` (default 1 hour). Before, one MFA kept the vault
and the admin settings unlocked for as long as the user stayed logged in.

This is a breaking change for deployments with `OCIS_MFA_ENABLED=true`. The
identity provider must:

* put the time of the MFA into the access token, by default as `auth_time`.
  Keycloak 25 and later does this by default. Use `OCIS_MFA_ISSUED_AT_CLAIM`
  for a different claim.
* ask for the second factor again once it is older than the MFA session
  duration. In Keycloak, set the max age of the MFA level (`loa-max-age`) to
  `OCIS_MFA_SESSION_DURATION` or lower.

Otherwise users can't open the vault or the admin settings after the first
hour. Set `OCIS_MFA_ENFORCE_SESSION_DURATION=false` to keep the old behaviour.

https://github.com/owncloud/ocis/pull/13538
