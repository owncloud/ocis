Enhancement: Update ownCloud Web to v6.0.0

Tags: web

We updated ownCloud Web to v6.0.0. Please refer to the changelog (linked) for details on the web release.

### Breaking changes
* BREAKING CHANGE for users in [owncloud/web#6648](https://github.com/owncloud/web/issues/6648): breaks existing bookmarks - they won't resolve anymore.
* BREAKING CHANGE for developers in [owncloud/web#6648](https://github.com/owncloud/web/issues/6648): the appDefaults composables from web-pkg now work with drive aliases, concatenated with relative item paths, instead of webdav paths. If you use the appDefaults composables in your application it's likely that your code needs to be adapted.

### Changes
* Bugfix [owncloud/web#7419](https://github.com/owncloud/web/issues/7419): Add language param opening external app
* Bugfix [owncloud/web#7731](https://github.com/owncloud/web/pull/7731): "Copy Quicklink"-translations
* Bugfix [owncloud/web#7830](https://github.com/owncloud/web/pull/7830): "Cut" and "Copy" actions for current folder
* Bugfix [owncloud/web#7652](https://github.com/owncloud/web/pull/7652): Disable copy/move overwrite on self
* Bugfix [owncloud/web#7739](https://github.com/owncloud/web/pull/7739): Disable shares loading on public and trash locations
* Bugfix [owncloud/web#7740](https://github.com/owncloud/web/pull/7740): Disappearing quicklink in sidebar
* Bugfix [owncloud/web#7946](https://github.com/owncloud/web/issues/7946): Prevent shares from disappearing after sharing with groups
* Bugfix [owncloud/web#7820](https://github.com/owncloud/web/pull/7820): Edit new created user in user management
* Bugfix [owncloud/web#7936](https://github.com/owncloud/web/pull/7936): Editing text files on public pages
* Bugfix [owncloud/web#7861](https://github.com/owncloud/web/pull/7861): Handle non 2xx external app responses
* Bugfix [owncloud/web#7734](https://github.com/owncloud/web/pull/7734): File name reactivity
* Bugfix [owncloud/web#7975](https://github.com/owncloud/web/pull/7975): Prevent file upload when folder creation failed
* Bugfix [owncloud/web#7724](https://github.com/owncloud/web/pull/7724): Folder conflict dialog
* Bugfix [owncloud/web#7603](https://github.com/owncloud/web/issues/7603): Hide search bar in public link context
* Bugfix [owncloud/web#7889](https://github.com/owncloud/web/pull/7889): Hide share indicators on public page
* Bugfix [owncloud/web#7903](https://github.com/owncloud/web/issues/7903): "Keep both"-conflict option
* Bugfix [owncloud/web#7697](https://github.com/owncloud/web/issues/7697): Link indicator on "Shared with me"-page
* Bugfix [owncloud/web#8007](https://github.com/owncloud/web/pull/8007): Missing password form on public drop page
* Bugfix [owncloud/web#7652](https://github.com/owncloud/web/pull/7652): Inhibit move files between spaces
* Bugfix [owncloud/web#7985](https://github.com/owncloud/web/pull/7985): Prevent retrying uploads with status code 5xx
* Bugfix [owncloud/web#7811](https://github.com/owncloud/web/pull/7811): Do not load files from cache in public links
* Bugfix [owncloud/web#7941](https://github.com/owncloud/web/pull/7941): Add origin check to Draw.io events
* Bugfix [owncloud/web#7916](https://github.com/owncloud/web/pull/7916): Prefer alias links over private links
* Bugfix [owncloud/web#7640](https://github.com/owncloud/web/pull/7640): "Private link"-button alignment
* Bugfix [owncloud/web#8006](https://github.com/owncloud/web/pull/8006): Public link loading on role change
* Bugfix [owncloud/web#7962](https://github.com/owncloud/web/issues/7962): Quota check when replacing files
* Bugfix [owncloud/web#7748](https://github.com/owncloud/web/pull/7748): Reload file list after last share removal
* Bugfix [owncloud/web#7699](https://github.com/owncloud/web/issues/7699): Remove the "close sidebar"-calls on delete
* Bugfix [owncloud/web#7504](https://github.com/owncloud/web/pull/7504): Resolve upload existing folder
* Bugfix [owncloud/web#7771](https://github.com/owncloud/web/pull/7771): Routing for re-shares
* Bugfix [owncloud/web#7675](https://github.com/owncloud/web/pull/7675): Search bar on small screens
* Bugfix [owncloud/web#7662](https://github.com/owncloud/web/pull/7662): Sidebar for received shares in search file list
* Bugfix [owncloud/web#7873](https://github.com/owncloud/web/pull/7873): Share editing after selecting a space
* Bugfix [owncloud/web#7657](https://github.com/owncloud/web/issues/7657): Share permissions for re-shares
* Bugfix [owncloud/web#7506](https://github.com/owncloud/web/issues/7506): Shares loading
* Bugfix [owncloud/web#7632](https://github.com/owncloud/web/pull/7632): Sidebar toggle icon
* Bugfix [owncloud/web#7781](https://github.com/owncloud/web/issues/7781): Sidebar without highlighted resource
* Bugfix [owncloud/web#7756](https://github.com/owncloud/web/pull/7756): Try to obtain refresh token before the error case
* Bugfix [owncloud/web#7768](https://github.com/owncloud/web/pull/7768): Hide actions in space trash bins
* Bugfix [owncloud/web#7651](https://github.com/owncloud/web/pull/7651): Spaces on "Shared via link"-page
* Bugfix [owncloud/web#7521](https://github.com/owncloud/web/issues/7521): Spaces reactivity on update
* Bugfix [owncloud/web#7960](https://github.com/owncloud/web/issues/7960): Display error messages in text editor
* Bugfix [owncloud/web#8030](https://github.com/owncloud/web/pull/8030): Saving a file multiple times with the text editor
* Bugfix [owncloud/web#7778](https://github.com/owncloud/web/issues/7778): Trash bin sidebar
* Bugfix [owncloud/web#7956](https://github.com/owncloud/web/issues/7956): Introduce "upload finalizing"-state in upload overlay
* Bugfix [owncloud/web#7630](https://github.com/owncloud/web/pull/7630): Upload modify time
* Bugfix [owncloud/web#8011](https://github.com/owncloud/web/issues/8011): Prevent unnecessary request when saving a user
* Bugfix [owncloud/web#7989](https://github.com/owncloud/web/pull/7989): Versions on the "Shared with me"-page
* Change [owncloud/web#6648](https://github.com/owncloud/web/issues/6648): Drive aliases in URLs
* Change [owncloud/web#7935](https://github.com/owncloud/web/pull/7935): Remove mediaSource and v-image-source
* Enhancement [owncloud/web#7635](https://github.com/owncloud/web/pull/7635): Add restore conflict dialog
* Enhancement [owncloud/web#7901](https://github.com/owncloud/web/pull/7901): Add search field for space members
* Enhancement [owncloud/web#4675](https://github.com/owncloud/web/issues/4675): Add `X-Request-ID` header to all outgoing requests
* Enhancement [owncloud/web#7904](https://github.com/owncloud/web/pull/7904): Batch actions for two or more items only
* Enhancement [owncloud/web#7892](https://github.com/owncloud/web/pull/7892): Respect the new sharing denials capability (experimental)
* Enhancement [owncloud/web#7709](https://github.com/owncloud/web/pull/7709): Edit custom permissions wording
* Enhancement [owncloud/web#7373](https://github.com/owncloud/web/issues/7373): Align dark mode colors with given design
* Enhancement [owncloud/web#7190](https://github.com/owncloud/web/pull/7190): Deny subfolders inside share
* Enhancement [owncloud/web#7684](https://github.com/owncloud/web/pull/7684): Design polishing
* Enhancement [owncloud/web#7865](https://github.com/owncloud/web/pull/7865): Disable share renaming
* Enhancement [owncloud/web#7725](https://github.com/owncloud/web/pull/7725): Enable renaming on received shares
* Enhancement [owncloud/web#7747](https://github.com/owncloud/web/pull/7747): Friendlier logout screen
* Enhancement [owncloud/web#6247](https://github.com/owncloud/web/issues/6247): Id based routing
* Enhancement [owncloud/web#7803](https://github.com/owncloud/web/issues/7803): Internal link on unaccepted share
* Enhancement [owncloud/web#7304](https://github.com/owncloud/web/issues/7304): Resolve internal links
* Enhancement [owncloud/web#7569](https://github.com/owncloud/web/pull/7569): Make keybindings global
* Enhancement [owncloud/web#7894](https://github.com/owncloud/web/pull/7894): Optimize email validation in the user management app
* Enhancement [owncloud/web#7707](https://github.com/owncloud/web/issues/7707): Resolve private links
* Enhancement [owncloud/web#7234](https://github.com/owncloud/web/issues/7234): Auth context in route meta props
* Enhancement [owncloud/web#7821](https://github.com/owncloud/web/pull/7821): Improve search experience
* Enhancement [owncloud/web#7801](https://github.com/owncloud/web/pull/7801): Make search results sortable
* Enhancement [owncloud/web#8028](https://github.com/owncloud/web/pull/8028): Update ODS to v14.0.1
* Enhancement [owncloud/web#7890](https://github.com/owncloud/web/pull/7890): Validate space names
* Enhancement [owncloud/web#7430](https://github.com/owncloud/web/pull/7430): Webdav support in web-client package
* Enhancement [owncloud/web#7900](https://github.com/owncloud/web/issues/7900): XHR upload timeout

https://github.com/owncloud/ocis/pull/5153
https://github.com/owncloud/web/releases/tag/v6.0.0
