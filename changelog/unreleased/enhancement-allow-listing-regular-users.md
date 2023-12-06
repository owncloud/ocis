Enhancement: Allow regular users to list other users

Regular users can search for other users. The following limitations
apply:

* Only search queries are allowed (using the `$search=term` query parameter)
* The search term needs to have at least 3 characters
* The result set only contains the attribute `displayName`, `userType`, `mail`
  and `id`

https://github.com/owncloud/ocis/pull/7887
https://github.com/owncloud/ocis/issues/7782
