Bugfix: Fix the `ocis search` command

We've fixed the behavior for `ocis search`, which didn't show further help when not all secrets have been configured.
It also was not possible to start the search service standalone from the oCIS binary without configuring all oCIS secrets,
even they were not needed by the search service.

https://github.com/owncloud/ocis/pull/3796

