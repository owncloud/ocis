Enhancement: Replace embedded IDP React SPA with server-rendered login page

Replaced the embedded IDP login's React SPA (pnpm, 21k LOC) with server-rendered HTML keeping theming and localiation support. The login page now works with minimal JavaScript, loads faster, and has a much smaller dependecy vulnerability surface.

https://github.com/owncloud/ocis/pull/12086
