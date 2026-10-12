Bugfix: Announce status messages to screen readers

Status messages such as "copied to clipboard" or "file restored successfully"
were not reliably read out by screen readers. The underlying announcement
area was never cleared between messages, so a repeated or rapid message
could be silently dropped, especially when combined with a focus change
elsewhere on the page.

The announcement area is now cleared before each new message and refilled
after a short delay, so status messages are announced consistently.

https://github.com/owncloud/ocis/pull/13024
