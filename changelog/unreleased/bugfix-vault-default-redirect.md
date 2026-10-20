Bugfix: Redirect to the default app when opening the vault root

Opening the vault without a path (`/vault`) showed a "not found" page, because
only the regular root redirected to the default app. The vault root now
redirects to the default app inside the vault, the same way the regular root
does outside of it.

https://github.com/owncloud/ocis/pull/13192
