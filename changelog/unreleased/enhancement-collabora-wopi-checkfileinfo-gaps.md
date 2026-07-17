Enhancement: Close Collabora WOPI CheckFileInfo property gaps

The collaboration service's Collabora CheckFileInfo response was missing several
properties that the Microsoft and OnlyOffice fileinfo types already had `case`
handling for, but that were never populated in the shared properties map:
Version, LastModifiedTime, ReadOnly, SupportsUpdate, IsAnonymousUser, several
file URLs (CloseUrl, DownloadUrl, plus the already-computed HostEditUrl/
HostViewUrl/FileSharingUrl/FileVersionUrl), the BreadcrumbBrandName/
BreadcrumbBrandUrl breadcrumb pair, an EditModePostMessage flag, and two new
Collabora-specific properties, IsUserLocked and IsAdminUser.

We've added all of these. Note that since Microsoft and OnlyOffice already had
`case` handling for ReadOnly and CloseUrl (and Microsoft for DownloadUrl too),
populating those keys unconditionally changes their CheckFileInfo output as
well, not just Collabora's: ReadOnly now actually surfaces as `true` for
read-only/view-only sessions on all three targets, and CloseUrl (plus
DownloadUrl for Microsoft) is now populated where it previously wasn't. Only
IsUserLocked and IsAdminUser are genuinely Collabora-only additions with no
matching `case` elsewhere.

https://github.com/owncloud/ocis/pull/12593
