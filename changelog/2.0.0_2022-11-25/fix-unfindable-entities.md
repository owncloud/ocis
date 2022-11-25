Bugfix: Fix unfindable entities from shares/publicshares

We fixed a problem where directories or empty files weren't findable because
they were to the search index improperly when created through a share or
publicshare.

https://github.com/owncloud/ocis/pull/4651
https://github.com/owncloud/ocis/issues/4489
