Enhancement: Pick the role of multi-instance guests from their claims

In a multi-instance setup every guest got the same role,
`OCIS_MULTI_INSTANCE_GUEST_ROLE`, no matter what their roles claim said. A
teacher invited to another instance could therefore not be given access to
vault (Safe) resources there without giving it to every other guest as well.

The new `OCIS_MULTI_INSTANCE_GUEST_ROLE_MAPPING` setting takes a list of
`claimvalue:rolename` pairs, for example `teacher:guest-teacher`. A guest whose
`PROXY_ROLE_ASSIGNMENT_OIDC_CLAIM` claim matches a pair gets that role instead
of the guest role, guests without a match keep `OCIS_MULTI_INSTANCE_GUEST_ROLE`.
A pair that points to a role which does not exist fails the login of the
matching guests. The proxy does not start if the setting is malformed or used
without the `oidc` role assignment driver. The setting is empty by default,
which keeps the previous behaviour.

https://github.com/owncloud/ocis/pull/13231
