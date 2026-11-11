Bugfix: Announce IDP sign-in errors and fix keyboard use

The login, consent, and account chooser pages did not announce error
messages when they appeared, and invalid username/password fields gave no
indication of their error state. The consent screen's scope list also had no
programmatic link between each entry and its checkbox, and the account
chooser's "continue as" and "use another account" entries could not be
reached or activated with a keyboard. A visible focus outline on the login
form's input fields had also been removed with no replacement. All of these
have been fixed.

https://github.com/owncloud/ocis/pull/12649
