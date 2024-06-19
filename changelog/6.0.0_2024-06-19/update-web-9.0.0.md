Enhancement: Update web to v9.0.0-alpha.7

Tags: web

We updated ownCloud Web to v9.0.0-alpha.7. Please refer to the changelog (linked) for details on the web release.

* Bugfix [owncloud/web#10377](https://github.com/owncloud/web/pull/10377): User data not updated while altering own user
* Bugfix [owncloud/web#10417](https://github.com/owncloud/web/pull/10417): Admin settings keyboard navigation
* Bugfix [owncloud/web#10517](https://github.com/owncloud/web/pull/10517): Load thumbnail when postprocessing is finished
* Bugfix [owncloud/web#10551](https://github.com/owncloud/web/pull/10551): Share sidebar icons
* Bugfix [owncloud/web#10702](https://github.com/owncloud/web/pull/10702): Apply sandbox attribute to iframe in draw-io extension
* Bugfix [owncloud/web#10706](https://github.com/owncloud/web/pull/10706): Apply sandbox attribute to iframe in app-external extension
* Bugfix [owncloud/web#10746](https://github.com/owncloud/web/pull/10746): Versions loaded multiple times when opening sidebar
* Bugfix [owncloud/web#10760](https://github.com/owncloud/web/pull/10760): Incoming notifications broken while notification center is open
* Bugfix [owncloud/web#10814](https://github.com/owncloud/web/issues/10814): Vertical scroll for OcModal on small screens
* Bugfix [owncloud/web#10900](https://github.com/owncloud/web/pull/10900): Context menu empty in tiles view
* Bugfix [owncloud/web#10918](https://github.com/owncloud/web/issues/10918): Resource deselection on right-click
* Bugfix [owncloud/web#10920](https://github.com/owncloud/web/pull/10920): Resources with name consist of number won't show up in trash bin
* Bugfix [owncloud/web#10928](https://github.com/owncloud/web/pull/10928): Disable search in public link context
* Bugfix [owncloud/web#10941](https://github.com/owncloud/web/issues/10941): Space not updating on navigation
* Bugfix [owncloud/web#10974](https://github.com/owncloud/web/pull/10974): Local logout if IdP has no logout support
* Change [owncloud/web#7338](https://github.com/owncloud/web/issues/7338): Remove deprecated code
* Change [owncloud/web#9892](https://github.com/owncloud/web/issues/9892): Remove skeleton app
* Change [owncloud/web#10102](https://github.com/owncloud/web/pull/10102): Remove deprecated extension point for adding quick actions
* Change [owncloud/web#10122](https://github.com/owncloud/web/pull/10122): Remove homeFolder option
* Change [owncloud/web#10210](https://github.com/owncloud/web/issues/10210): Vuex store removed
* Change [owncloud/web#10240](https://github.com/owncloud/web/pull/10240): Remove ocs user
* Change [owncloud/web#10330](https://github.com/owncloud/web/pull/10330): Registering app file editors
* Change [owncloud/web#10443](https://github.com/owncloud/web/pull/10443): Add extensionPoint concept
* Change [owncloud/web#10758](https://github.com/owncloud/web/pull/10758): Portal target removed
* Change [owncloud/web#10786](https://github.com/owncloud/web/pull/10786): Disable opening files in embed mode
* Enhancement [owncloud/web#5383](https://github.com/owncloud/web/issues/5383): Accessibility improvements
* Enhancement [owncloud/web#9215](https://github.com/owncloud/web/issues/9215): Icon for .dcm files
* Enhancement [owncloud/web#10018](https://github.com/owncloud/web/issues/10018): Tile sizes
* Enhancement [owncloud/web#10207](https://github.com/owncloud/web/pull/10207): Enable user preferences in public links
* Enhancement [owncloud/web#10334](https://github.com/owncloud/web/pull/10334): Move ThemeSwitcher into Account Settings
* Enhancement [owncloud/web#10383](https://github.com/owncloud/web/issues/10383): Top loading bar increase visibility
* Enhancement [owncloud/web#10390](https://github.com/owncloud/web/pull/10390): Integrate ToastUI editor in the text editor app
* Enhancement [owncloud/web#10443](https://github.com/owncloud/web/pull/10443): Custom component extension type
* Enhancement [owncloud/web#10448](https://github.com/owncloud/web/pull/10448): Epub reader app
* Enhancement [owncloud/web#10485](https://github.com/owncloud/web/pull/10485): Highlight search term in sharing autosuggest list
* Enhancement [owncloud/web#10519](https://github.com/owncloud/web/pull/10519): Warn user before closing browser when upload is in progress
* Enhancement [owncloud/web#10534](https://github.com/owncloud/web/issues/10534): Full text search default
* Enhancement [owncloud/web#10544](https://github.com/owncloud/web/pull/10544): Show locked and processing next to other status indicators
* Enhancement [owncloud/web#10546](https://github.com/owncloud/web/pull/10546): Set emoji as space icon
* Enhancement [owncloud/web#10586](https://github.com/owncloud/web/pull/10586): Add SSE events for locking, renaming, deleting, and restoring
* Enhancement [owncloud/web#10611](https://github.com/owncloud/web/pull/10611): Remember left nav bar state
* Enhancement [owncloud/web#10612](https://github.com/owncloud/web/pull/10612): Remember right side bar state
* Enhancement [owncloud/web#10624](https://github.com/owncloud/web/pull/10624): Add details panel to trash
* Enhancement [owncloud/web#10709](https://github.com/owncloud/web/pull/10709): Implement Server-Sent Events (SSE) for File Creation
* Enhancement [owncloud/web#10758](https://github.com/owncloud/web/pull/10758): Search providers extension point
* Enhancement [owncloud/web#10782](https://github.com/owncloud/web/pull/10782): Implement Server-Sent Events (SSE) for file updates
* Enhancement [owncloud/web#10798](https://github.com/owncloud/web/pull/10798): Add SSE event for moving
* Enhancement [owncloud/web#10801](https://github.com/owncloud/web/pull/10801): Ability to theme sharing role icons
* Enhancement [owncloud/web#10807](https://github.com/owncloud/web/pull/10807): Add SSE event for moving
* Enhancement [owncloud/web#10874](https://github.com/owncloud/web/pull/10874): Show loading spinner while searching or filtering users
* Enhancement [owncloud/web#10907](https://github.com/owncloud/web/pull/10907): Display hidden resources information in files list
* Enhancement [owncloud/web#10929](https://github.com/owncloud/web/pull/10929): Add loading spinner to admin settings spaces and groups
* Enhancement [owncloud/web#10956](https://github.com/owncloud/web/pull/10956): Audio metadata panel
* Enhancement [owncloud/web#10956](https://github.com/owncloud/web/pull/10956): EXIF metadata panel
* Enhancement [owncloud/web#10976](https://github.com/owncloud/web/pull/10976): Faster page loading times
* Enhancement [owncloud/web#11004](https://github.com/owncloud/web/pull/11004): Add enabled only filter to spaces overview
* Enhancement [owncloud/web#11037](https://github.com/owncloud/web/pull/11037): Multiple sidebar root panels

https://github.com/owncloud/ocis/pull/9395
https://github.com/owncloud/web/releases/tag/v9.0.0
