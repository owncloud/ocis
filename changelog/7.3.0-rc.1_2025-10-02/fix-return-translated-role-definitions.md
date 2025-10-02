Bugfix: Return translated role definitions

Instead of always returning the role definitions in English, we now return the role definitions in the language set in the `Accept-Language` header if present.

https://github.com/owncloud/ocis/pull/11466
