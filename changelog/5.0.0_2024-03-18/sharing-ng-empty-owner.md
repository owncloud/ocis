Bugfix: graph/sharedWithMe works for shares from project spaces now

We fixed a bug in the 'graph/v1beta1/me/drive/sharedWithMe' endpoint that
caused an error response when the user received shares from project spaces.
Additionally the endpoint now behaves more graceful in cases where the
displayname of the owner or creator of a share or shared resource couldn't be
resolved.

https://github.com/owncloud/ocis/pull/8233
https://github.com/owncloud/ocis/issues/8027
https://github.com/owncloud/ocis/issues/8215
