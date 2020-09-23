Change: Tighten screws on usernames and email addresses

In order to match accounts to the OIDC claims we currently rely on the email address or username to be present. We force both to match the [W3C recommended regex](https://www.w3.org/TR/2016/REC-html51-20161101/sec-forms.html#valid-e-mail-address) with usernames having to start with a character or `_`. This allows the username to be presented and used in ACLs when integrating the os with the glauth LDAP service of ocis.

<https://github.com/owncloud/ocis/accounts/pull/65>
