Bugfix: Show the Safe access-denied page in the user's language

A user without permission to use the Safe who opened a Safe URL directly saw
the access-denied page in the browser language instead of the language set in
their account. It only switched to the account language after clicking
"Continue" and going back. The page now opens in the account language.

https://github.com/owncloud/ocis/pull/13500
