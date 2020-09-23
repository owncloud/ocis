Bugfix: Fix loading and saving system scoped values

We fixed loading and saving system scoped values. Those are now saved without an account uuid, so that the value
can be loaded by other accounts as well.

<https://github.com/owncloud/ocis/settings/pull/66>
