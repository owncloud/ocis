Bugfix: Recover the login page after the current user was deleted

When a user account was deleted while its session was still stored in the
browser, opening the web interface redirected to an authentication error page
and looped there. The identity provider kept accepting the stale logon session
but could no longer resolve the deleted user, so the sign-in flow ended up
re-issuing the authorization request without its parameters and the browser was
sent back to the same error over and over. The only way out was to clear the
browser data manually before a different user could sign in.

The identity provider now detects this stuck authorization loop, expires the
stale logon session on the server side and restarts a clean login, so the sign-in
form is shown again and a different user can sign in right away.

https://github.com/owncloud/ocis/pull/12851
