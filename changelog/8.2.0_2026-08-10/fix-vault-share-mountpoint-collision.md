Bugfix: Do not collide vault and non-vault share mountpoints

Received shares mounted from the Vault storage and received shares mounted
from regular drives were treated as one flat namespace when computing a
unique mountpoint name. Sharing a resource with the same name once from a
regular drive and once from the Vault caused the second one to get a
spurious ` (1)` suffix, even though the two shares are rendered in
completely separate, segregated lists and never actually collide.
Mountpoint name collisions are now only checked against other shares within
the same vault/non-vault group.

https://github.com/owncloud/ocis/pull/12730
