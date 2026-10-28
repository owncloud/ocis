Enhancement: Allow group shares in the vault

Users can now invite a group to a resource in the vault (Safe). Before this
change, the server accepted only users as vault share recipients. The share
dialog in Web now also searches for groups in vault mode. The
notifications and userlog services now send a vault notification only to the
group members who have the permission to use the vault. The other group
members get no email and no in-app notification.

https://github.com/owncloud/ocis/pull/13259
https://github.com/owncloud/reva/pull/787
