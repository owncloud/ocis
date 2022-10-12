Enhancement: Update ownCloud Web to v6.0.0-rc.2

Tags: web

We updated ownCloud Web to v6.0.0-rc.2. Please refer to the changelog (linked) for details on the web release.

### Breaking changes
* BREAKING CHANGE for users in [owncloud/web#6648](https://github.com/owncloud/web/issues/6648): breaks existing bookmarks - they won't resolve anymore.
* BREAKING CHANGE for developers in [owncloud/web#6648](https://github.com/owncloud/web/issues/6648): the appDefaults composables from web-pkg now work with drive aliases, concatenated with relative item paths, instead of webdav paths. If you use the appDefaults composables in your application it's likely that your code needs to be adapted.

### Changes
* Bugfix [owncloud/web#7419](https://github.com/owncloud/web/issues/7419): Add language param opening external app
* Bugfix [owncloud/web#7731](https://github.com/owncloud/web/pull/7731): "Copy Quicklink"-translations
* Bugfix [owncloud/web#7652](https://github.com/owncloud/web/pull/7652): Disable copy/move overwrite on self
* Bugfix [owncloud/web#7739](https://github.com/owncloud/web/pull/7739): Disable shares loading on public and trash locations
* Bugfix [owncloud/web#7740](https://github.com/owncloud/web/pull/7740): Disappearing quicklink in sidebar
* Bugfix [owncloud/web#7734](https://github.com/owncloud/web/pull/7734): File name reactivity
* Bugfix [owncloud/web#7724](https://github.com/owncloud/web/pull/7724): Folder conflict dialog
* Bugfix [owncloud/web#7652](https://github.com/owncloud/web/pull/7652): Inhibit move files between spaces
* Bugfix [owncloud/web#7640](https://github.com/owncloud/web/pull/7640): "Private link"-button alignment
* Bugfix [owncloud/web#7748](https://github.com/owncloud/web/pull/7748): Reload file list after last share removal
* Bugfix [owncloud/web#7699](https://github.com/owncloud/web/issues/7699): Remove the "close sidebar"-calls on delete
* Bugfix [owncloud/web#7504](https://github.com/owncloud/web/pull/7504): Resolve upload existing folder
* Bugfix [owncloud/web#7675](https://github.com/owncloud/web/pull/7675): Search bar on small screens
* Bugfix [owncloud/web#7662](https://github.com/owncloud/web/pull/7662): Sidebar for received shares in search file list
* Bugfix [owncloud/web#7506](https://github.com/owncloud/web/issues/7506): Shares loading
* Bugfix [owncloud/web#7632](https://github.com/owncloud/web/pull/7632): Sidebar toggle icon
* Bugfix [owncloud/web#7756](https://github.com/owncloud/web/pull/7756): Try to obtain refresh token before the error case
* Bugfix [owncloud/web#7651](https://github.com/owncloud/web/pull/7651): Spaces on "Shared via link"-page
* Bugfix [owncloud/web#7521](https://github.com/owncloud/web/issues/7521): Spaces reactivity on update
* Bugfix [owncloud/web#7630](https://github.com/owncloud/web/pull/7630): Upload modify time
* Change [owncloud/web#6648](https://github.com/owncloud/web/issues/6648): Drive aliases in URLs
* Enhancement [owncloud/web#7709](https://github.com/owncloud/web/pull/7709): Edit custom permissions wording
* Enhancement [owncloud/web#7190](https://github.com/owncloud/web/pull/7190): Deny subfolders inside share
* Enhancement [owncloud/web#7684](https://github.com/owncloud/web/pull/7684): Design polishing
* Enhancement [owncloud/web#7725](https://github.com/owncloud/web/pull/7725): Enable renaming on received shares
* Enhancement [owncloud/web#7747](https://github.com/owncloud/web/pull/7747): Friendlier logout screen
* Enhancement [owncloud/web#6247](https://github.com/owncloud/web/issues/6247): Id based routing
* Enhancement [owncloud/web#7405](https://github.com/owncloud/web/pull/7405): Resolve internal links
* Enhancement [owncloud/web#7569](https://github.com/owncloud/web/pull/7569): Make keybindings global
* Enhancement [owncloud/web#7405](https://github.com/owncloud/web/pull/7405): Resolve private links
* Enhancement [owncloud/web#7684](https://github.com/owncloud/web/pull/7684): Update ODS to v14.0.0-alpha.20
* Enhancement [owncloud/web#7430](https://github.com/owncloud/web/pull/7430): Webdav support in web-client package

https://github.com/owncloud/ocis/pull/4786
https://github.com/owncloud/web/releases/tag/v6.0.0-rc.2
