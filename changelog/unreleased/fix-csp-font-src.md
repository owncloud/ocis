Bugfix: Fix CSP blocking bundled KaTeX font

The default Content Security Policy blocked the bundled KaTeX math font (used by the md-editor) because it is inlined as a `data:` URI in the Web UI CSS. Added `data:` to the `font-src` directive to resolve the console error on every page load. Users with custom CSP files (`PROXY_CSP_CONFIG_FILE_LOCATION`) will need to add `data:` to their `font-src` directive manually.

https://github.com/owncloud/ocis/pull/12070
