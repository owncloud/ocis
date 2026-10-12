Bugfix: Show correct message when Safe access is denied for lacking permission

A user without permission to use the Safe who opened the Safe URL directly was
shown "Not logged in" with a "Log in again" button, even though they were
still logged in. The page now shows "Access denied" with a "Go back" button
that returns into Drive for this case.

https://github.com/owncloud/ocis/pull/13111
