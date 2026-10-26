Enhancement: Pick the role of multi-instance guests from their claims

In a multi-instance setup every guest got the same role,
`OCIS_MULTI_INSTANCE_GUEST_ROLE`, no matter what their roles claim said. A
teacher invited to another instance could therefore not be given access to
vault (Safe) resources there without giving it to every other guest as well.

The new `OCIS_MULTI_INSTANCE_GUEST_ROLE_MAPPING` setting takes a list of
`claimvalue:rolename` pairs, for example `teacher:guest-teacher`. When the
`oidc` role assignment driver is used, a guest whose
`PROXY_ROLE_ASSIGNMENT_OIDC_CLAIM` claim matches a pair gets that role instead
of the guest role. Guests without a match, or whose mapped role does not exist,
keep `OCIS_MULTI_INSTANCE_GUEST_ROLE`. The setting is empty by default, which
keeps the previous behaviour.

https://github.com/owncloud/ocis/pull/13231
