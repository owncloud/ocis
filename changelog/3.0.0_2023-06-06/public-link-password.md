Enhancement: Add config option to enforce passwords on public links

Added a new config option to enforce passwords on public links with "Uploader, Editor, Contributor" roles.

The new options are:
`OCIS_SHARING_PUBLIC_WRITEABLE_SHARE_MUST_HAVE_PASSWORD`, `SHARING_PUBLIC_WRITEABLE_SHARE_MUST_HAVE_PASSWORD` and `FRONTEND_OCS_PUBLIC_WRITEABLE_SHARE_MUST_HAVE_PASSWORD`.
Check the docs on how to properly set them.

https://github.com/owncloud/ocis/pull/5848
https://github.com/owncloud/ocis/pull/5785
https://github.com/owncloud/ocis/pull/5720
