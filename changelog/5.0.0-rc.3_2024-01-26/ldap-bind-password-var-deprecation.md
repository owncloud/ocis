Enhancement: Add OCIS_LDAP_BIND_PASSWORD as replacement for LDAP_BIND_PASSWORD

The enviroment variable `OCIS_LDAP_BIND_PASSWORD` was added to be more consistent with all
other global LDAP variables.

`LDAP_BIND_PASSWORD` is deprecated now and scheduled for removal with the 5.0.0 release.

We also deprecated `LDAP_USER_SCHEMA_ID_IS_OCTETSTRING` for removal with 5.0.0.
The replacement for it is `OCIS_LDAP_USER_SCHEMA_ID_IS_OCTETSTRING`.

https://github.com/owncloud/ocis/issues/7176
