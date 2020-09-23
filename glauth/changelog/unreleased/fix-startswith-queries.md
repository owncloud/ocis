Bugfix: fix LDAP substring startswith filters

Filters like `(mail=mar*)` are currentld not parsed correctly, but they are used when searching for recipients. This PR correctly converts them to odata filters like `startswith(mail,'mar')`.

<https://github.com/owncloud/ocis/glauth/pull/31>
