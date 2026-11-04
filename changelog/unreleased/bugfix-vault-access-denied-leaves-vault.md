Bugfix: Leave the Safe when access to it is denied

A user without permission to use the Safe who opened a Safe URL directly was
denied access, but the page stayed in Safe mode. Clicking "Continue" then
crashed the application, and the browser back button showed the user's regular
files with the Safe look. The denial is now shown outside the Safe, so every
following navigation returns to the regular app.

https://github.com/owncloud/ocis/pull/13343
