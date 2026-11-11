Bugfix: Hide app menu items from anonymous public link visitors

The app drawer in the top bar offered "Activities" and "Library" to anonymous
visitors following a public link. Both apps registered their app menu item
unconditionally at bootstrap, before any authentication context exists, and the
extension registry is shared between the authenticated and the public link
context. Because the public link file view renders the full application layout,
the top bar and its app drawer came along with it - even though the sibling top
bar elements (notification bell, sidebar toggle, home link) already gate
themselves on the authentication context.

Following the leaked entries did not expose any data: both apps route with
`authContext: 'user'`, so the visitor was redirected to the login page. What
leaked was the fact that these apps are enabled on the instance.

Both apps now contribute their menu item only when a user is signed in, matching
what the files, text editor, app store, admin settings and ScienceMesh apps
already do. With no menu items left to show, the existing top bar condition
hides the app drawer on public links by itself.

https://github.com/owncloud/ocis/pull/12842

