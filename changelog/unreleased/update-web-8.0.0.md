Enhancement: Update web to v8.0.0

Tags: web

We updated ownCloud Web to v8.0.0. Please refer to the changelog (linked) for details on the web release.

* Bugfix [owncloud/web#9257](https://github.com/owncloud/web/issues/9257): Filter out shares without display name
* Bugfix [owncloud/web#9529](https://github.com/owncloud/web/pull/9529): Shared with action menu label alignment
* Bugfix [owncloud/web#9649](https://github.com/owncloud/web/pull/9649): Add project space filter
* Bugfix [owncloud/web#9663](https://github.com/owncloud/web/pull/9663): Respect the open-in-new-tab-config for external apps
* Bugfix [owncloud/web#9694](https://github.com/owncloud/web/issues/9694): Special characters in username
* Bugfix [owncloud/web#9788](https://github.com/owncloud/web/issues/9788): Create .space folder if it does not exist
* Bugfix [owncloud/web#9799](https://github.com/owncloud/web/issues/9799): Link resolving into default app
* Bugfix [owncloud/web#9832](https://github.com/owncloud/web/pull/9832): Copy quicklinks for webkit navigator
* Bugfix [owncloud/web#9843](https://github.com/owncloud/web/pull/9843): Fix display path on resources
* Bugfix [owncloud/web#9844](https://github.com/owncloud/web/pull/9844): Upload space image
* Bugfix [owncloud/web#9861](https://github.com/owncloud/web/pull/9861): Duplicated file search request
* Bugfix [owncloud/web#9873](https://github.com/owncloud/web/pull/9873): Tags are no longer editable for a locked file
* Bugfix [owncloud/web#9881](https://github.com/owncloud/web/pull/9881): Prevent rendering of old/wrong set of resources in search list
* Bugfix [owncloud/web#9915](https://github.com/owncloud/web/pull/9915): Keep both folders conflict in same-named folders
* Bugfix [owncloud/web#9931](https://github.com/owncloud/web/pull/9931): Enabling "invite people" for password-protected folder/file
* Bugfix [owncloud/web#10010](https://github.com/owncloud/web/issues/10010): Displaying full video in their dimensions
* Bugfix [owncloud/web#10031](https://github.com/owncloud/web/issues/10031): Icon extension mapping
* Bugfix [owncloud/web#10065](https://github.com/owncloud/web/pull/10065): Logout page after token expiry
* Bugfix [owncloud/web#10083](https://github.com/owncloud/web/pull/10083): Disable expiration date for alias link (internal)
* Bugfix [owncloud/web#10092](https://github.com/owncloud/web/pull/10092): Allow empty search query in "in-here" search
* Bugfix [owncloud/web#10096](https://github.com/owncloud/web/pull/10096): Remove password buttons on input if disabled
* Bugfix [owncloud/web#10118](https://github.com/owncloud/web/pull/10118): Tilesview has whitespace
* Bugfix [owncloud/web#10149](https://github.com/owncloud/web/pull/10149): Spaces files list previews cropped
* Bugfix [owncloud/web#10149](https://github.com/owncloud/web/pull/10149): Spaces overview tile previews zoomed
* Bugfix [owncloud/web#10154](https://github.com/owncloud/web/pull/10154): Resolving links without drive alias
* Bugfix [owncloud/web#10156](https://github.com/owncloud/web/pull/10156): Uploading the same files parallel
* Bugfix [owncloud/web#10158](https://github.com/owncloud/web/pull/10158): GDPR export polling
* Bugfix [owncloud/web#10176](https://github.com/owncloud/web/pull/10176): Turned off file extensions not always respected
* Bugfix [owncloud/web#10179](https://github.com/owncloud/web/pull/10179): Space navigate to trash missing
* Bugfix [owncloud/web#10182](https://github.com/owncloud/web/pull/10182): Make versions panel readonly in viewers and editors
* Bugfix [owncloud/web#10220](https://github.com/owncloud/web/pull/10220): Loading indicator during conflict dialog
* Bugfix [owncloud/web#10227](https://github.com/owncloud/web/issues/10227): Configurable concurrent requests
* Bugfix [owncloud/web#10232](https://github.com/owncloud/web/pull/10232): Skip searchbar preview fetch on reload
* Bugfix [owncloud/web#10318](https://github.com/owncloud/web/pull/10318): Scrollable account page
* Bugfix [owncloud/web#10321](https://github.com/owncloud/web/pull/10321): Private link error messages
* Bugfix [owncloud/web#10347](https://github.com/owncloud/web/pull/10347): Readonly user attributes have no effect on group memberships
* Bugfix [owncloud/web#10424](https://github.com/owncloud/web/pull/10424): Restore space
* Bugfix [owncloud/web#10473](https://github.com/owncloud/web/issues/10473): Public link file download
* Bugfix [owncloud/web#10489](https://github.com/owncloud/web/pull/10489): Wrong share permissions when resharing off
* Bugfix [owncloud/web#10514](https://github.com/owncloud/web/pull/10514): Indicate shares that are not manageable due to file locking
* Change [owncloud/web#2404](https://github.com/owncloud/web/issues/2404): Theme handling
* Change [owncloud/web#7338](https://github.com/owncloud/web/issues/7338): Remove deprecated code
* Change [owncloud/web#9653](https://github.com/owncloud/web/pull/9653): Keyword Query Language (KQL) search syntax
* Change [owncloud/web#9709](https://github.com/owncloud/web/issues/9709): DavProperties without namespace
* Enhancement [owncloud/web#7317](https://github.com/owncloud/ocis/pull/7317): Make login url configurable
* Enhancement [owncloud/web#7497](https://github.com/owncloud/ocis/issues/7497): Permission checks for shares and favorites
* Enhancement [owncloud/web#7600](https://github.com/owncloud/web/issues/7600): Scroll to newly created folder
* Enhancement [owncloud/web#9302](https://github.com/owncloud/web/issues/9302): Application unification
* Enhancement [owncloud/web#9423](https://github.com/owncloud/web/pull/9423): Show local loading spinner in sharing button
* Enhancement [owncloud/web#9441](https://github.com/owncloud/web/pull/9441): File versions tooltip with absolute date
* Enhancement [owncloud/web#9441](https://github.com/owncloud/web/pull/9441): Disabling extensions
* Enhancement [owncloud/web#9451](https://github.com/owncloud/web/pull/9451): Add SSE to get notifications instantly
* Enhancement [owncloud/web#9525](https://github.com/owncloud/web/pull/9525): Tags form improved
* Enhancement [owncloud/web#9527](https://github.com/owncloud/web/pull/9527): Don't display confirmation dialog on file deletion
* Enhancement [owncloud/web#9531](https://github.com/owncloud/web/issues/9531): Personal shares can be shown and hidden
* Enhancement [owncloud/web#9552](https://github.com/owncloud/web/pull/9552): Upload preparation time
* Enhancement [owncloud/web#9561](https://github.com/owncloud/web/pull/9561): Indicate processing state
* Enhancement [owncloud/web#9566](https://github.com/owncloud/web/pull/9566): Display locking information
* Enhancement [owncloud/web#9584](https://github.com/owncloud/web/pull/9584): Moving share's "set expiration date" function
* Enhancement [owncloud/web#9625](https://github.com/owncloud/web/pull/9625): Add keyboard navigation to spaces overview
* Enhancement [owncloud/web#9627](https://github.com/owncloud/web/pull/9627): Add batch actions to spaces
* Enhancement [owncloud/web#9671](https://github.com/owncloud/web/pull/9671): OcModal set buttons to same width
* Enhancement [owncloud/web#9682](https://github.com/owncloud/web/pull/9682): Add password policy compatibility
* Enhancement [owncloud/web#9691](https://github.com/owncloud/web/pull/9691): Password generator for public links
* Enhancement [owncloud/web#9696](https://github.com/owncloud/web/pull/9696): Added app banner for mobile devices
* Enhancement [owncloud/web#9706](https://github.com/owncloud/web/pull/9706): Unify sharing expiration date menu items
* Enhancement [owncloud/web#9709](https://github.com/owncloud/web/issues/9709): New WebDAV implementation in web-client
* Enhancement [owncloud/web#9727](https://github.com/owncloud/web/pull/9727): Show error if password is on a banned password list
* Enhancement [owncloud/web#9768](https://github.com/owncloud/web/issues/9768): Embed mode
* Enhancement [owncloud/web#9771](https://github.com/owncloud/web/pull/9771): Handle postprocessing state via Server Sent Events
* Enhancement [owncloud/web#9794](https://github.com/owncloud/web/pull/9794): Registering search providers as extension
* Enhancement [owncloud/web#9806](https://github.com/owncloud/web/pull/9806): Preview image presentation
* Enhancement [owncloud/web#9809](https://github.com/owncloud/web/pull/9809): Add editors to the application menu
* Enhancement [owncloud/web#9814](https://github.com/owncloud/web/pull/9814): Registering nav items as extension
* Enhancement [owncloud/web#9815](https://github.com/owncloud/web/pull/9815): Add new portal into runtime to include footer
* Enhancement [owncloud/web#9831](https://github.com/owncloud/web/pull/9831): Last modified filter chips
* Enhancement [owncloud/web#9847](https://github.com/owncloud/web/issues/9847): Provide vendor neutral file icons
* Enhancement [owncloud/web#9854](https://github.com/owncloud/web/pull/9854): Search query term linking
* Enhancement [owncloud/web#9857](https://github.com/owncloud/web/pull/9857): Add permission to delete link passwords when password is enforced
* Enhancement [owncloud/web#9858](https://github.com/owncloud/web/pull/9858): Remove settings icon from searchbar
* Enhancement [owncloud/web#9864](https://github.com/owncloud/web/pull/9864): Search tags filter chips style aligned
* Enhancement [owncloud/web#9884](https://github.com/owncloud/web/pull/9884): Enable dark theme on importer
* Enhancement [owncloud/web#9890](https://github.com/owncloud/web/pull/9890): Create shortcuts
* Enhancement [owncloud/web#9905](https://github.com/owncloud/web/pull/9905): Manage tags in details panel
* Enhancement [owncloud/web#9906](https://github.com/owncloud/web/pull/9906): Reorganize "New" menu
* Enhancement [owncloud/web#9912](https://github.com/owncloud/web/pull/9912): Add media type filter chip
* Enhancement [owncloud/web#9940](https://github.com/owncloud/web/pull/9940): Display error message for upload to locked folder
* Enhancement [owncloud/web#9966](https://github.com/owncloud/web/issues/9966): Support more audio formats with correct icon
* Enhancement [owncloud/web#10007](https://github.com/owncloud/web/issues/10007): Additional languages
* Enhancement [owncloud/web#10013](https://github.com/owncloud/web/issues/10013): Shared by filter
* Enhancement [owncloud/web#10014](https://github.com/owncloud/web/issues/10014): Share search filter
* Enhancement [owncloud/web#10024](https://github.com/owncloud/web/pull/10024): Duplicate space
* Enhancement [owncloud/web#10037](https://github.com/owncloud/web/pull/10037): Default link permission
* Enhancement [owncloud/web#10047](https://github.com/owncloud/web/pull/10047): Add explaining contextual helper to spaces overview
* Enhancement [owncloud/web#10057](https://github.com/owncloud/web/pull/10057): Folder tree creation during upload
* Enhancement [owncloud/web#10062](https://github.com/owncloud/web/pull/10062): Show webdav information in details view
* Enhancement [owncloud/web#10099](https://github.com/owncloud/web/pull/10099): Support mandatory filter while listing users
* Enhancement [owncloud/web#10102](https://github.com/owncloud/web/pull/10102): Registering quick actions as extension
* Enhancement [owncloud/web#10104](https://github.com/owncloud/web/pull/10104): Create link modal
* Enhancement [owncloud/web#10111](https://github.com/owncloud/web/pull/10111): Registering right sidebar panels as extension
* Enhancement [owncloud/web#10111](https://github.com/owncloud/web/pull/10111): File sidebar in viewer and editor apps
* Enhancement [owncloud/web#10224](https://github.com/owncloud/web/pull/10224): Harmonize AppSwitcher icon colors
* Enhancement [owncloud/web#10356](https://github.com/owncloud/web/pull/10356): Preview app add reset button for images

https://github.com/owncloud/ocis/pull/8613
https://github.com/owncloud/web/releases/tag/v8.0.0
