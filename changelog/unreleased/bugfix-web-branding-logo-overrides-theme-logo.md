Bugfix: Admin-uploaded logo now overrides theme-specific logos

Previously, a theme's own per-variant logo defined under
`clients.web.themes[].logo` would keep showing even after an admin uploaded a
custom logo via the admin settings web app, because the upload only patched
the `clients.web.defaults.logo` keys. This affected the built-in "owncloud"
theme too, since most of its variants (dark, high contrast, ...) already
define their own logo. The admin-uploaded logo now always overrides a
theme's own per-variant logo.

https://github.com/owncloud/ocis/pull/12847
