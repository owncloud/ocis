Bugfix: Keep the upload details list scrollable

After uploading or extracting many files at once, expanding "Show details" in
the upload panel grew the list past the top of the screen. The list could not
be scrolled and "Hide details" was out of reach until the page was reloaded.
The panel now stays within the screen and only the list scrolls, so the
title bar and buttons stay visible, also with failed uploads and on small screens.

https://github.com/owncloud/ocis/pull/13494
