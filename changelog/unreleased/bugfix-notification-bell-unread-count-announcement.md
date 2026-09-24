Bugfix: notification bell now announces the unread count

The notification bell's accessible name was a static "Notifications"
string regardless of how many unread notifications there were. Since
that label overrides all other content for assistive technology, screen
reader users never heard the unread count that sighted users see on the
bell's badge.

The bell's accessible name and tooltip now include the exact unread
count, with correct singular/plural wording.

https://github.com/owncloud/ocis/pull/13003
