Bugfix: Mark all notifications as read

Clicking "Mark all read" in the notifications menu only cleared the list
visually. The request sent to the server was malformed, so the notifications
were not deleted and reappeared after a page reload. This affected the regular
view as well as the Safe.

The request is now sent correctly and the notifications are removed for good.

https://github.com/owncloud/ocis/pull/13157
