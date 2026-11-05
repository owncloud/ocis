Bugfix: Hide the space "Edit description" button without edit permission

The pencil button next to the description in the space header was shown to
every space member, including those who are not allowed to edit the
description. It is now only shown to users who are allowed to edit the space
description, consistent with the equivalent entries in the space context menu
and the sidebar.

https://github.com/owncloud/ocis/pull/13350
