Enhancement: add `ocis.id` and numeric id claims

We added an `ocis.id` claim to the OIDC standard claims. It allows the idp to send a stable identifier that can be exposed to the outside world (in contrast to sub, which might change whens the IdP changes).

In addition we added `uidnumber` and `gidnumber` claims, which can be used by the IdP as well. They will be used by storage providers that integrate with an existing LDAP server.

https://github.com/owncloud/ocis-pkg/pull/50
