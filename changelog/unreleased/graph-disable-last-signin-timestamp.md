Enhancement: Allow disabling the last sign-in timestamp update

The graph service maintains the 'oCLastSignInTimestamp' LDAP attribute of a user
on every sign-in (when the LDAP identity backend has write access). This can
cause a significant amount of LDAP write load, especially when the proxy's OIDC
userinfo cache has a short TTL and sign-in events are emitted frequently.

A new setting 'OCIS_LDAP_UPDATE_LAST_SIGNIN_DATE' / 'GRAPH_LDAP_UPDATE_LAST_SIGNIN_DATE'
(default 'true') allows disabling the update of the last sign-in timestamp
without having to disable all LDAP writes ('OCIS_LDAP_SERVER_WRITE_ENABLED') or
the graph events consumer. When set to 'false' the graph service no longer
listens for 'UserSignedIn' events and does not write the 'oCLastSignInTimestamp'
attribute.

https://github.com/owncloud/ocis/pull/12522
https://github.com/owncloud/ocis/issues/9942
