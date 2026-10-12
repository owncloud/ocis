Bugfix: fixed arrow-key navigation in the notification dropdown

Opening the notification bell's dropdown did not move keyboard focus into
it. Arrow-key navigation only started working after an extra, unrelated
Tab press, and pressing an arrow key too early scrolled the page instead,
which looked like the dropdown had closed.

The dropdown now moves focus to its first item as soon as it opens, so
arrow keys work immediately, matching the behavior of other click-triggered
dropdowns in the application.

https://github.com/owncloud/ocis/pull/13002
