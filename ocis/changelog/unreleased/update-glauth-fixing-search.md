Bugfix: Update ocis-glauth for fixed single user search

We updated ocis-glauth to a version that comes with a fix for searching a single user or group. ocis-glauth was dropping search context before by ignoring the searchBaseDN for filtering. This has been fixed.

<https://github.com/owncloud/product/issues/214>
<https://github.com/owncloud/ocis/pull/535>
<https://github.com/owncloud/ocis-glauth/pull/32>
