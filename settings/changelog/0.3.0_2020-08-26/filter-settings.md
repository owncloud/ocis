Change: Filter settings by permissions

`BundleService.GetBundle` and `BundleService.ListBundles` are now filtered by READ permissions in the role of the authenticated user. This prevents settings from being visible to the user when their role doesn't have appropriate permissions.

<https://github.com/owncloud/product/issues/99>
<https://github.com/owncloud/ocis/settings/pull/48>
