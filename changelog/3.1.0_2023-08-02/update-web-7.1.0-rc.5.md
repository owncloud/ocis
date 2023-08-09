Enhancement: Update web to v7.1.0-rc.5

Tags: web

We updated ownCloud Web to v7.1.0-rc.5. Please refer to the changelog (linked) for details on the web release.

## Summary
* Bugfix [owncloud/web#9078](https://github.com/owncloud/web/pull/9078): Favorites list update on removal
* Bugfix [owncloud/web#9213](https://github.com/owncloud/web/pull/9213): Space creation does not block reoccurring event
* Bugfix [owncloud/web#9247](https://github.com/owncloud/web/issues/9247): Uploading to folders that contain special characters
* Bugfix [owncloud/web#9259](https://github.com/owncloud/web/issues/9259): Relative user quota display limited to two decimals
* Bugfix [owncloud/web#9261](https://github.com/owncloud/web/issues/9261): Remember location after token invalidation
* Bugfix [owncloud/web#9299](https://github.com/owncloud/web/pull/9299): Authenticated public links breaking uploads
* Bugfix [owncloud/web#9315](https://github.com/owncloud/web/issues/9315): Switch columns displayed on small screens in "Shared with me" view
* Bugfix [owncloud/web#9351](https://github.com/owncloud/web/pull/9351): Media controls overflow on mobile screens
* Bugfix [owncloud/web#9389](https://github.com/owncloud/web/pull/9389): Space editors see empty trashbin and delete actions in space trashbin
* Bugfix [owncloud/web#9461](https://github.com/owncloud/web/pull/9461): Merging folders
* Bugfix [owncloud/web/#9496](https://github.com/owncloud/web/pull/9496): Logo not showing
* Bugfix [owncloud/web/#9489](https://github.com/owncloud/web/pull/9489): Public drop zone
* Bugfix [owncloud/web/#9487](https://github.com/owncloud/web/pull/9487): Respect supportedClouds config
* Bugfix [owncloud/web/#9507](https://github.com/owncloud/web/pull/9507): Space description edit modal is cut off vertically
* Bugfix [owncloud/web/#9501](https://github.com/owncloud/web/pull/9501): Add cloud importer translations
* Bugfix [owncloud/web/#9510](https://github.com/owncloud/web/pull/9510): Double items after moving a file with the same name
* Enhancement [owncloud/web#7967](https://github.com/owncloud/web/pull/7967): Add hasPriority property for editors per extension
* Enhancement [owncloud/web#8422](https://github.com/owncloud/web/issues/8422): Improve extension app topbar
* Enhancement [owncloud/web#8445](https://github.com/owncloud/web/issues/8445): Open individually shared file in dedicated view
* Enhancement [owncloud/web#8599](https://github.com/owncloud/web/issues/8599): Shrink table columns
* Enhancement [owncloud/web#8921](https://github.com/owncloud/web/pull/8921): Add whitespace context-menu
* Enhancement [owncloud/web#8983](https://github.com/owncloud/web/pull/8983): Deny share access
* Enhancement [owncloud/web#8984](https://github.com/owncloud/web/pull/8984): Long breadcrumb strategy
* Enhancement [owncloud/web#9044](https://github.com/owncloud/web/pull/9044): Search tag filter
* Enhancement [owncloud/web#9046](https://github.com/owncloud/web/pull/9046): Single file link open with default app
* Enhancement [owncloud/web#9052](https://github.com/owncloud/web/pull/9052): Drag & drop on parent folder
* Enhancement [owncloud/web#9055](https://github.com/owncloud/web/pull/9055): Respect archiver limits
* Enhancement [owncloud/web#9056](https://github.com/owncloud/web/issues/9056): Enable download (archive) on spaces
* Enhancement [owncloud/web#9059](https://github.com/owncloud/web/pull/9059): Search full-text filter
* Enhancement [owncloud/web#9077](https://github.com/owncloud/web/pull/9077): Advanced search button
* Enhancement [owncloud/web#9077](https://github.com/owncloud/web/pull/9077): Search breadcrumb
* Enhancement [owncloud/web#9088](https://github.com/owncloud/web/pull/9088): Use app icons for files
* Enhancement [owncloud/web#9140](https://github.com/owncloud/web/pull/9140): Upload file on paste
* Enhancement [owncloud/web#9151](https://github.com/owncloud/web/issues/9151): Cloud import
* Enhancement [owncloud/web#9174](https://github.com/owncloud/web/issues/9174): Privacy statement in account menu
* Enhancement [owncloud/web#9178](https://github.com/owncloud/web/pull/9178): Add login button to top bar
* Enhancement [owncloud/web#9195](https://github.com/owncloud/web/pull/9195): Project spaces list viewmode
* Enhancement [owncloud/web#9199](https://github.com/owncloud/web/pull/9199): Add pagination options to admin settings
* Enhancement [owncloud/web#9200](https://github.com/owncloud/web/pull/9200): Add batch actions to search result list
* Enhancement [owncloud/web#9216](https://github.com/owncloud/web/issues/9216): Restyle possible sharees
* Enhancement [owncloud/web#9226](https://github.com/owncloud/web/pull/9226): Streamline URL query names
* Enhancement [owncloud/web#9263](https://github.com/owncloud/web/pull/9263): Access denied page update message
* Enhancement [owncloud/web#9280](https://github.com/owncloud/web/issues/9280): Hover tooltips in topbar
* Enhancement [owncloud/web#9294](https://github.com/owncloud/web/pull/9294): Search list add highlighted file content
* Enhancement [owncloud/web#9299](https://github.com/owncloud/web/pull/9299): Resolve pulic links to their actual location
* Enhancement [owncloud/web#9304](https://github.com/owncloud/web/pull/9304): Add search location filter
* Enhancement [owncloud/web#9344](https://github.com/owncloud/web/pull/9344): Ambiguation for URL view mode params
* Enhancement [owncloud/web#9346](https://github.com/owncloud/web/pull/9346): Batch actions redesign
* Enhancement [owncloud/web#9348](https://github.com/owncloud/web/pull/9348): Tag comma separation on client side
* Enhancement [owncloud/web#9377](https://github.com/owncloud/web/issues/9377): User notification for blocked pop-ups and redirects
* Enhancement [owncloud/web#9386](https://github.com/owncloud/web/pull/9386): Allow local storage for auth token
* Enhancement [owncloud/web#9394](https://github.com/owncloud/web/pull/9394): Button styling
* Enhancement [owncloud/web#9449](https://github.com/owncloud/web/issues/9449): Error notifications include x-request-id
* Enhancement [owncloud/web#9426](https://github.com/owncloud/web/pull/9426): Add error log to upload dialog


https://github.com/owncloud/ocis/pull/6944
https://github.com/owncloud/web/releases/tag/v7.1.0-rc.5
