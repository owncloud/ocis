Enhancement: Validate space names

We now return `BAD REQUEST` when space names are
-  too long (max 255 characters)
-  containing evil characters (`/`, `\`, `.`, `\\`, `:`, `?`, `*`, `"`, `>`, `<`, `|`)

Additionally leading and trailing spaces will be removed silently.

https://github.com/owncloud/ocis/pull/4955
