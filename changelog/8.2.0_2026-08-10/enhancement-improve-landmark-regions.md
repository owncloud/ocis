Enhancement: Improve page structure on several pages

We've added missing `main` and `footer` regions to several pages (404,
private link resolving, missing config, logout, access denied, public link
resolving) for a more consistent page structure. We've also fixed the
"Skip to main" link, which previously did nothing on pages using the plain
layout (login, logout, public/private link pages) due to a missing target id
and a stale cached DOM reference.

https://github.com/owncloud/ocis/pull/12637
