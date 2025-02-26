Bugfix: apply role constraints when creating shares via the graph API

We fixed a bug in the graph API for creating and updating shares so that
Spaceroot specific roles like 'Manager' and 'Co-owner' can no longer be
assigned for shares on files or directories.

https://github.com/owncloud/ocis/pull/8247
https://github.com/owncloud/ocis/issues/8131
