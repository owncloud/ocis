Enhancement: Support Apache Tika 4.x in content extraction

The search service can now extract content from an Apache Tika 4.x server in
addition to Tika 3.x. Tika 4.0 renamed the `X-TIKA:content` metadata key to
`tk:content`, replaced the `/language/string` endpoint with `/language` and
reports the GPS altitude as `geo:alt` instead of `GPS:GPS Altitude`. Without
support for these changes, a 4.x server would return no content for oCIS, so file
contents were no longer searchable, stop word cleaning was skipped and the
altitude of photos was missing.

The extractor now reads `tk:content` and falls back to `X-TIKA:content`, reads
the altitude from both `geo:alt` and `GPS:GPS Altitude`, and picks the language
endpoint from the version the Tika server reports at startup. Both Tika 3.x and
4.x servers are supported.

https://github.com/owncloud/ocis/pull/13279
