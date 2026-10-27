Enhancement: Make the user search attributes configurable

Which attributes user search matches against is now configurable via
`OCIS_LDAP_USER_SEARCH_ATTRIBUTES`, defaulting to name, display name and
e-mail. Deployments that need stricter matching (e.g. for data protection)
can restrict it to name, surname and username, excluding e-mail.

https://github.com/owncloud/ocis/pull/13219
