Enhancement: Add option to disable public link sharing

Added an `OCIS_ENABLE_PUBLIC_SHARING` config option, read by the frontend
and sharing services. It defaults to `true`. When set to `false`,
creating new public links is rejected and the
`files_sharing.public.enabled` capability reports `false` so clients
hide the corresponding UI. Direct sharing with users and groups is not
affected.

https://github.com/owncloud/ocis/pull/12542
