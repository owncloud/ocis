Enhancement: Update web to v7.0.0-rc.37

Tags: web

We updated ownCloud Web to v7.0.0-rc.37. Please refer to the changelog (linked) for details on the web release.

* Bugfix [owncloud/web#6423](https://github.com/owncloud/web/issues/6423): Archiver in protected public links
* Bugfix [owncloud/web#6434](https://github.com/owncloud/web/issues/6434): Endless lazy loading indicator after sorting file table
* Bugfix [owncloud/web#6731](https://github.com/owncloud/web/issues/6731): Layout with long breadcrumb
* Bugfix [owncloud/web#6768](https://github.com/owncloud/web/issues/6768): Pagination after increasing items per page
* Bugfix [owncloud/web#7513](https://github.com/owncloud/web/issues/7513): Calendar popup position in right sidebar
* Bugfix [owncloud/web#7655](https://github.com/owncloud/web/issues/7655): Loading shares in deep nested folders
* Bugfix [owncloud/web#7925](https://github.com/owncloud/web/pull/7925): "Paste"-action without write permissions
* Bugfix [owncloud/web#7926](https://github.com/owncloud/web/pull/7926): Include spaces in the list info
* Bugfix [owncloud/web#7958](https://github.com/owncloud/web/pull/7958): Prevent deletion of own account
* Bugfix [owncloud/web#7966](https://github.com/owncloud/web/pull/7966): UI fixes for sorting and quickactions
* Bugfix [owncloud/web#7969](https://github.com/owncloud/web/pull/7969): Space quota not displayed after creation
* Bugfix [owncloud/web#8026](https://github.com/owncloud/web/pull/8026): Text editor appearance
* Bugfix [owncloud/web#8040](https://github.com/owncloud/web/pull/8040): Reverting versions for read-only shares
* Bugfix [owncloud/web#8045](https://github.com/owncloud/web/pull/8045): Resolving drives in search
* Bugfix [owncloud/web#8054](https://github.com/owncloud/web/issues/8054): Search repeating no results message
* Bugfix [owncloud/web#8058](https://github.com/owncloud/web/pull/8058): Current year selection in the date picker
* Bugfix [owncloud/web#8061](https://github.com/owncloud/web/pull/8061): Omit "page"-query in breadcrumb navigation
* Bugfix [owncloud/web#8080](https://github.com/owncloud/web/pull/8080): Left sidebar navigation item text flickers on transition
* Bugfix [owncloud/web#8081](https://github.com/owncloud/web/issues/8081): Space member disappearing
* Bugfix [owncloud/web#8083](https://github.com/owncloud/web/issues/8083): Re-using space images
* Bugfix [owncloud/web#8148](https://github.com/owncloud/web/issues/8148): Show space members despite deleted entries
* Bugfix [owncloud/web#8158](https://github.com/owncloud/web/issues/8158): Search bar input appearance
* Bugfix [owncloud/web#8265](https://github.com/owncloud/web/pull/8265): Application menu active display on hover
* Bugfix [owncloud/web#8276](https://github.com/owncloud/web/pull/8276): Loading additional user data
* Bugfix [owncloud/web#8300](https://github.com/owncloud/web/pull/8300): Re-loading space members panel
* Bugfix [owncloud/web#8326](https://github.com/owncloud/web/pull/8326): Editing users who never logged in
* Bugfix [owncloud/web#8340](https://github.com/owncloud/web/pull/8340): Cancel custom permissions
* Bugfix [owncloud/web#8411](https://github.com/owncloud/web/issues/8411): Drop menus with limited vertical screen space
* Bugfix [owncloud/web#8420](https://github.com/owncloud/web/issues/8420): Token renewal in vue router hash mode
* Bugfix [owncloud/web#8434](https://github.com/owncloud/web/issues/8434): Accessing route in admin-settings with insufficient permissions
* Bugfix [owncloud/web#8479](https://github.com/owncloud/web/issues/8479): "Show more"-action in shares panel
* Bugfix [owncloud/web#8480](https://github.com/owncloud/web/pull/8480): Paste action conflict dialog broken
* Bugfix [owncloud/web#8498](https://github.com/owncloud/web/pull/8498): PDF display issue - Update CSP object-src policy
* Bugfix [owncloud/web#8508](https://github.com/owncloud/web/pull/8508): Remove fuzzy search results
* Bugfix [owncloud/web#8523](https://github.com/owncloud/web/issues/8523): Space image upload
* Bugfix [owncloud/web#8549](https://github.com/owncloud/web/issues/8549): Batch context actions in admin settings
* Bugfix [owncloud/web#8554](https://github.com/owncloud/web/pull/8554): Height of dropdown no-option
* Bugfix [owncloud/web#8576](https://github.com/owncloud/web/pull/8576): De-duplicate event handling to prevent errors on Draw-io
* Bugfix [owncloud/web#8585](https://github.com/owncloud/web/issues/8585): Users without role assignment
* Bugfix [owncloud/web#8587](https://github.com/owncloud/web/issues/8587): Password enforced check for public links
* Bugfix [owncloud/web#8592](https://github.com/owncloud/web/issues/8592): Group members sorting
* Bugfix [owncloud/web#8694](https://github.com/owncloud/web/pull/8694): Broken re-login after logout
* Bugfix [owncloud/web#8695](https://github.com/owncloud/web/issues/8695): Open files in external app
* Bugfix [owncloud/web#8756](https://github.com/owncloud/web/pull/8756): Copy link to clipboard text
* Bugfix [owncloud/web#8758](https://github.com/owncloud/web/pull/8758): Preview controls colors
* Bugfix [owncloud/web#8776](https://github.com/owncloud/web/issues/8776): Selection reset on action click
* Bugfix [owncloud/web#8814](https://github.com/owncloud/web/pull/8814): Share recipient container exceed
* Bugfix [owncloud/web#8825](https://github.com/owncloud/web/pull/8825): Remove drop target in read-only folders
* Bugfix [owncloud/web#8827](https://github.com/owncloud/web/pull/8827): Opening context menu via keyboard
* Bugfix [owncloud/web#8834](https://github.com/owncloud/web/issues/8834): Hide upload hint in empty read-only folders
* Bugfix [owncloud/web#8864](https://github.com/owncloud/web/pull/8864): Public link empty password stays forever
* Bugfix [owncloud/web#8880](https://github.com/owncloud/web/issues/8880): Sidebar header after deleting resource
* Bugfix [owncloud/web#8928](https://github.com/owncloud/web/issues/8928): Infinite login redirect
* Bugfix [owncloud/web#8987](https://github.com/owncloud/web/pull/8987): Limit amount of concurrent tus requests
* Bugfix [owncloud/web#8992](https://github.com/owncloud/web/pull/8992): Personal space name after language change
* Bugfix [owncloud/web#9004](https://github.com/owncloud/web/issues/9004): Endless loading when encountering a public link error
* Bugfix [owncloud/web#9015](https://github.com/owncloud/web/pull/9015): Prevent "virtual" spaces from being displayed in the UI
* Change [owncloud/web#6661](https://github.com/owncloud/web/issues/6661): Streamline new tab handling in extensions
* Change [owncloud/web#7948](https://github.com/owncloud/web/issues/7948): Update Vue to v3.2
* Change [owncloud/web#8431](https://github.com/owncloud/web/pull/8431): Remove permission manager
* Change [owncloud/web#8455](https://github.com/owncloud/web/pull/8455): Configurable extension autosave
* Change [owncloud/web#8563](https://github.com/owncloud/web/pull/8563): Theme colors
* Enhancement [owncloud/web#6183](https://github.com/owncloud/web/issues/6183): Global loading indicator
* Enhancement [owncloud/web#7388](https://github.com/owncloud/web/pull/7388): Add tag support
* Enhancement [owncloud/web#7721](https://github.com/owncloud/web/issues/7721): Improve performance when loading folders and share indicators
* Enhancement [owncloud/web#7942](https://github.com/owncloud/web/pull/7942): Warn users when using unsupported browsers
* Enhancement [owncloud/web#7965](https://github.com/owncloud/web/pull/7965): Optional Contributor role and configurable resharing permissions
* Enhancement [owncloud/web#7968](https://github.com/owncloud/web/pull/7968): Group and user creation forms submit on enter
* Enhancement [owncloud/web#7976](https://github.com/owncloud/web/pull/7976): Add switch to enable condensed resource table
* Enhancement [owncloud/web#7977](https://github.com/owncloud/web/pull/7977): Introduce zoom and rotate to the preview app
* Enhancement [owncloud/web#7983](https://github.com/owncloud/web/pull/7983): Conflict dialog UX
* Enhancement [owncloud/web#7991](https://github.com/owncloud/web/pull/7991): Add tiles view for resource display
* Enhancement [owncloud/web#7994](https://github.com/owncloud/web/pull/7994): Introduce full screen mode to the preview app
* Enhancement [owncloud/web#7995](https://github.com/owncloud/web/pull/7995): Enable autoplay in the preview app
* Enhancement [owncloud/web#8008](https://github.com/owncloud/web/issues/8008): Don't open sidebar when copying quicklink
* Enhancement [owncloud/web#8021](https://github.com/owncloud/web/pull/8021): Access right sidebar panels via URL
* Enhancement [owncloud/web#8051](https://github.com/owncloud/web/pull/8051): Introduce image preloading to the preview app
* Enhancement [owncloud/web#8055](https://github.com/owncloud/web/pull/8055): Retry failed uploads on re-upload
* Enhancement [owncloud/web#8056](https://github.com/owncloud/web/pull/8056): Increase Searchbar height
* Enhancement [owncloud/web#8057](https://github.com/owncloud/web/pull/8057): Show text file icon for empty text files
* Enhancement [owncloud/web#8132](https://github.com/owncloud/web/pull/8132): Update libre-graph-api to v1.0
* Enhancement [owncloud/web#8136](https://github.com/owncloud/web/pull/8136): Make clipboard copy available to more browsers
* Enhancement [owncloud/web#8161](https://github.com/owncloud/web/pull/8161): Space group members
* Enhancement [owncloud/web#8161](https://github.com/owncloud/web/pull/8161): Space group shares
* Enhancement [owncloud/web#8166](https://github.com/owncloud/web/issues/8166): Show upload speed
* Enhancement [owncloud/web#8175](https://github.com/owncloud/web/pull/8175): Rename "user management" app
* Enhancement [owncloud/web#8178](https://github.com/owncloud/web/pull/8178): Spaces list in admin settings
* Enhancement [owncloud/web#8261](https://github.com/owncloud/web/pull/8261): Admin settings users section uses graph api for role assignments
* Enhancement [owncloud/web#8279](https://github.com/owncloud/web/pull/8279): Move user group select to edit panel
* Enhancement [owncloud/web#8280](https://github.com/owncloud/web/pull/8280): Add support for multiple clients in `theme.json`
* Enhancement [owncloud/web#8294](https://github.com/owncloud/web/pull/8294): Move language selection to user account page
* Enhancement [owncloud/web#8306](https://github.com/owncloud/web/pull/8306): Show selectable groups only
* Enhancement [owncloud/web#8317](https://github.com/owncloud/web/pull/8317): Add context menu to groups
* Enhancement [owncloud/web#8320](https://github.com/owncloud/web/pull/8320): Space member expiration
* Enhancement [owncloud/web#8320](https://github.com/owncloud/web/pull/8320): Update SDK to v3.1.0-alpha.3
* Enhancement [owncloud/web#8324](https://github.com/owncloud/web/pull/8324): Add context menu to users
* Enhancement [owncloud/web#8331](https://github.com/owncloud/web/pull/8331): Admin settings users section details improvement
* Enhancement [owncloud/web#8354](https://github.com/owncloud/web/issues/8354): Add `ItemFilter` component
* Enhancement [owncloud/web#8356](https://github.com/owncloud/web/pull/8356): Slight improvement of key up/down performance
* Enhancement [owncloud/web#8363](https://github.com/owncloud/web/issues/8363): Admin settings general section
* Enhancement [owncloud/web#8375](https://github.com/owncloud/web/pull/8375): Add appearance section in general settings
* Enhancement [owncloud/web#8377](https://github.com/owncloud/web/issues/8377): User group filter
* Enhancement [owncloud/web#8387](https://github.com/owncloud/web/pull/8387): Batch edit quota in admin panel
* Enhancement [owncloud/web#8398](https://github.com/owncloud/web/pull/8398): Use standardized layout for file/space action list
* Enhancement [owncloud/web#8425](https://github.com/owncloud/web/issues/8425): Add dark ownCloud logo
* Enhancement [owncloud/web#8432](https://github.com/owncloud/web/pull/8432): Inject customizations
* Enhancement [owncloud/web#8433](https://github.com/owncloud/web/pull/8433): User settings login field
* Enhancement [owncloud/web#8441](https://github.com/owncloud/web/pull/8441): Skeleton App
* Enhancement [owncloud/web#8449](https://github.com/owncloud/web/pull/8449): Configurable top bar
* Enhancement [owncloud/web#8450](https://github.com/owncloud/web/pull/8450): Rework notification bell
* Enhancement [owncloud/web#8455](https://github.com/owncloud/web/pull/8455): Autosave content changes in text editor
* Enhancement [owncloud/web#8473](https://github.com/owncloud/web/pull/8473): Update CERN links
* Enhancement [owncloud/web#8489](https://github.com/owncloud/web/pull/8489): Respect max quota
* Enhancement [owncloud/web#8492](https://github.com/owncloud/web/pull/8492): User role filter
* Enhancement [owncloud/web#8503](https://github.com/owncloud/web/issues/8503): Beautify file version list
* Enhancement [owncloud/web#8515](https://github.com/owncloud/web/pull/8515): Introduce trashbin overview
* Enhancement [owncloud/web#8518](https://github.com/owncloud/web/pull/8518): Make notifications work with oCIS
* Enhancement [owncloud/web#8541](https://github.com/owncloud/web/pull/8541): Public link permission `PublicLink.Write.all`
* Enhancement [owncloud/web#8553](https://github.com/owncloud/web/pull/8553): Add and remove users from groups batch actions
* Enhancement [owncloud/web#8554](https://github.com/owncloud/web/pull/8554): Beautify form inputs
* Enhancement [owncloud/web#8557](https://github.com/owncloud/web/issues/8557): Rework mobile navigation
* Enhancement [owncloud/web#8566](https://github.com/owncloud/web/pull/8566): QuickActions role configurable
* Enhancement [owncloud/web#8612](https://github.com/owncloud/web/issues/8612): Add `Accept-Language` header to all outgoing requests
* Enhancement [owncloud/web#8630](https://github.com/owncloud/web/pull/8630): Add logout url
* Enhancement [owncloud/web#8652](https://github.com/owncloud/web/pull/8652): Enable guest users
* Enhancement [owncloud/web#8711](https://github.com/owncloud/web/pull/8711): Remove placeholder, add customizable label
* Enhancement [owncloud/web#8713](https://github.com/owncloud/web/pull/8713): Context helper read more link configurable
* Enhancement [owncloud/web#8715](https://github.com/owncloud/web/pull/8715): Enable rename groups
* Enhancement [owncloud/web#8730](https://github.com/owncloud/web/pull/8730): Create Space from selection
* Enhancement [owncloud/web#8738](https://github.com/owncloud/web/issues/8738): GDPR export
* Enhancement [owncloud/web#8762](https://github.com/owncloud/web/pull/8762): Stop bootstrapping application earlier in anonymous contexts
* Enhancement [owncloud/web#8766](https://github.com/owncloud/web/pull/8766): Add support for read-only groups
* Enhancement [owncloud/web#8790](https://github.com/owncloud/web/pull/8790): Custom translations
* Enhancement [owncloud/web#8797](https://github.com/owncloud/web/pull/8797): Font family in theming
* Enhancement [owncloud/web#8806](https://github.com/owncloud/web/pull/8806): Preview app sorting
* Enhancement [owncloud/web#8820](https://github.com/owncloud/web/pull/8820): Adjust missing reshare permissions message
* Enhancement [owncloud/web#8822](https://github.com/owncloud/web/pull/8822): Fix quicklink icon alignment
* Enhancement [owncloud/web#8826](https://github.com/owncloud/web/pull/8826): Admin settings groups members panel
* Enhancement [owncloud/web#8868](https://github.com/owncloud/web/pull/8868): Respect user read-only configuration by the server
* Enhancement [owncloud/web#8876](https://github.com/owncloud/web/pull/8876): Update roles and permissions names, labels, texts and icons
* Enhancement [owncloud/web#8882](https://github.com/owncloud/web/pull/8882): Layout of Share role and expiration date dropdown
* Enhancement [owncloud/web#8883](https://github.com/owncloud/web/issues/8883): Webfinger redirect app
* Enhancement [owncloud/web#8898](https://github.com/owncloud/web/pull/8898): Rename "Quicklink" to "link"
* Enhancement [owncloud/web#8911](https://github.com/owncloud/web/pull/8911): Add notification setting to account page

https://github.com/owncloud/ocis/pull/6294
https://github.com/owncloud/web/releases/tag/v7.0.0-rc.37
