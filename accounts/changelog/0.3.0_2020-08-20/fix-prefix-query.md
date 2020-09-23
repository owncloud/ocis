Bugfix: Unescape value for prefix query

Prefix queries also need to unescape token values like `'some ''ol string'` to `some 'ol string` before using it in a prefix query

<https://github.com/owncloud/ocis/accounts/pull/76>
