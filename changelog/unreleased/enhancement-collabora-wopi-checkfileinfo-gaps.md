Enhancement: Close Collabora WOPI CheckFileInfo property gaps

The collaboration service's Collabora CheckFileInfo response was missing several
properties that the Microsoft and OnlyOffice responses already had: Version,
LastModifiedTime, ReadOnly, SupportsUpdate, IsAnonymousUser, several file URLs
(CloseUrl, DownloadUrl, plus the already-computed HostEditUrl/HostViewUrl/
FileSharingUrl/FileVersionUrl), the BreadcrumbBrandName/BreadcrumbBrandUrl
breadcrumb pair, an EditModePostMessage flag, and two new Collabora-specific
properties, IsUserLocked and IsAdminUser.

We've added all of these.

https://github.com/owncloud/ocis/pull/12593
